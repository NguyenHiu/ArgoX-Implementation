package matcher

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (m *Matcher) Register() {
	m.prepareNonceAndGasPrice(1, 300000)

	_, err := m.OnchainInstance.Register(m.Auth, m.Address)
	if err != nil {
		log.Fatal(err)
	}

	go m.ListenEvents()
}

func (m *Matcher) ListenEvents() {
	opts := bind.WatchOpts{Context: context.Background()}
	go watchFullfilEvent(m.OnchainInstance, &opts)
}

func watchFullfilEvent(contract *onchain.Onchain, opts *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainFullfilMatch)
	sub, err := contract.WatchFullfilMatch(opts, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
			// case vLogs := <-logs:
			// 	_logger.Info("[Fullfill] batch id: %v\n", vLogs.Raw.Data)
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
