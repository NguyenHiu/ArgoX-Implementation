package matcher

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
				_orderID, _ := order.Data.OrderID.MarshalBinary()
				var _orderID16 [16]byte
				copy(_orderID16[:], _orderID)
				onchainOrders = append(onchainOrders, onchain.OnchainOrder{
					OrderID:   _orderID16,
					Price:     order.Data.Price,
					Amount:    order.Data.Amount,
					Side:      order.Data.Side,
					Signature: util.CorrectSignToOnchain(order.Data.Signature),
					Owner:     common.Address(common.FromHex(order.Data.Owner.String())),
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
	m.Auth.Value = client.EthToWei(big.NewFloat(float64(value)))
	m.Auth.GasLimit = uint64(gasLimit)
}
