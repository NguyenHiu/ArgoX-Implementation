package supermatcher

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/tradeClient"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (sm *SuperMatcher) SendBatch(batch *Batch) {
	sm.Mutex.Lock()

	_side := "bid"
	if batch.Side == constants.ASK {
		_side = "ask"
	}
	//IMHERETODEBUG_logger.Info("Sends batch::%v::%v::%v::%v to onchain\n", batch.BatchID.String()[:5], batch.Owner.String()[:5], batch.Amount, _side)

	sm.prepareNonceAndGasPrice(0, 900000)
	_signature := util.CorrectSignToOnchain(batch.Signature)
	_, err := sm.OnchainInstance.SendBatch(sm.Auth, batch.BatchID, batch.Price, batch.Amount, batch.Side, batch.Owner, _signature)

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
	sm.Auth.Value = tradeClient.EthToWei(big.NewFloat(float64(value)))
	sm.Auth.GasLimit = uint64(gasLimit)
}
