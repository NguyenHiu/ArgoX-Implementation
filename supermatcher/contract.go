package supermatcher

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (sm *SuperMatcher) SendBatch(batch *Batch) {
	sm.Mutex.Lock()

	_logger.Info("Sends batch::%v to onchain\n", batch.BatchID.String()[:6])

	sm.prepareNonceAndGasPrice(0, 900000)
	_, err := sm.OnchainInstance.SendBatch(sm.Auth, batch.BatchID, batch.Price, batch.Amount, batch.Side, batch.Owner, batch.Signature)

	sm.Mutex.Unlock()

	if err != nil {
		log.Fatal(err)
	}
}

func (sm *SuperMatcher) isMatcher(matcherAddr common.Address) bool {
	res, err := sm.OnchainInstance.IsMatcher(&bind.CallOpts{Context: context.Background()}, matcherAddr)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (sm *SuperMatcher) prepareNonceAndGasPrice(value float64, gasLimit int) {
	nonce, err := sm.Client.PendingNonceAt(context.Background(), sm.Address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := sm.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sm.Auth.Nonce = big.NewInt(int64(nonce))
	sm.Auth.GasPrice = gasPrice
	sm.Auth.Value = client.EthToWei(big.NewFloat(float64(value)))
	sm.Auth.GasLimit = uint64(gasLimit)
}
