package matcher

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/orderClient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

func (m *Matcher) Register() {
	m.prepareNonceAndGasPrice(1, 300000)

	_, err := m.OnchainInstance.Register(m.Auth, m.Address)
	if err != nil {
		log.Fatal(err)
	}

	m.ListenEvents()
}

func (m *Matcher) ListenEvents() {
	opts := bind.WatchOpts{Context: context.Background()}
	go m.watchFullfilEvent(&opts)
}

func (m *Matcher) watchFullfilEvent(opts *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainFullfilMatch)
	sub, err := m.OnchainInstance.WatchFullfilMatch(opts, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			// _logger.Debug("Matcher::%v receive an fullfill event batch::%v\n", m.ID.String()[:6], id.String())
			m.Mux.Lock()
			batch, ok := m.Batches[id]
			m.Mux.Unlock()
			if !ok {
				continue
			}

			onchainOrders := []onchain.OnchainOrder{}
			for _, order := range batch.Orders {
				trades := new(bytes.Buffer)
				for _, trade := range order.Trades {
					tData, err := trade.Encode_TransferBatching()
					if err != nil {
						_logger.Error("sending batch's detail got error, encode trade, err: %v\n", err)
					}
					if err := binary.Write(trades, binary.BigEndian, tData); err != nil {
						_logger.Error("sending batch's detail, hash a trade got error, err: %v\n", err)
					}
				}
				tradeHash := crypto.Keccak256Hash(trades.Bytes())

				oData, err := order.OriginalOrder.Encode_TransferBatching()
				if err != nil {
					_logger.Error("sending batch's detail, hash original order got error, err: %v\n", err)
				}
				originalOrderHash := crypto.Keccak256Hash(oData)

				onchainOrders = append(onchainOrders, onchain.OnchainOrder{
					Price:             order.ShadowOrder.Price,
					Amount:            order.ShadowOrder.Amount,
					Side:              order.ShadowOrder.Side,
					From:              order.ShadowOrder.From,
					TradeHash:         tradeHash,
					OriginalOrderHash: originalOrderHash,
					Owner:             common.Address(*order.OriginalOrder.Owner),
				})
			}

			// Send batch's details
			// m.Mux.Lock()
			m.prepareNonceAndGasPrice(0, 900000)
			_, err := m.OnchainInstance.SubmitOrderDetails(m.Auth, vLogs.Arg0, onchainOrders)
			if err != nil {
				_logger.Error("Submit Error, err: %v\n", err)
			}
			_logger.Debug("Matcher::%v, Submit batch's details, batch::%v\n", m.ID.String()[:6], id.String())
			// m.Mux.Unlock()
		}
	}
}

func (m *Matcher) prepareNonceAndGasPrice(value float64, gasLimit int) {
	nonce, err := m.Client.PendingNonceAt(context.Background(), m.Address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	m.Auth.Nonce = big.NewInt(int64(nonce))
	m.Auth.GasPrice = gasPrice
	m.Auth.Value = orderClient.EthToWei(big.NewFloat(float64(value)))
	m.Auth.GasLimit = uint64(gasLimit)
}
