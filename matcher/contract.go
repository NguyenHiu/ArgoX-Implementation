package matcher

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/orderClient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
			m.Mux.Lock()
			batch, ok := m.Batches[id]
			m.Mux.Unlock()
			if !ok {
				continue
			}

			m.BatchStatusMapping[id] = WAITING_PROOF

			if batch.Side == constants.BID {
				_logger.Debug("Matched, amount: %v\n", batch.Amount)
			}
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
