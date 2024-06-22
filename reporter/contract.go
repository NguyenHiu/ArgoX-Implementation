package reporter

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/google/uuid"
)

func (r *Reporter) Reporting() {
	go r.reporting()
}

func (r *Reporter) reporting() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		r.Mux.Lock()
		if len(r.PendingBatches) == 0 {
			r.Mux.Unlock()
			continue
		}

		for k, v := range r.PendingBatches {
			if time.Now().Unix()-v >= r.WaitingTime {
				_logger.Info("Report batch::%v\n", k.String())
				r.prepareNonceAndGasPrice(0, 1000000)
				batchID, _ := k.MarshalBinary()
				var _batchID [16]byte
				copy(_batchID[:], batchID[:16])
				_, err := r.OnchainInstance.ReportMissingDeadline(r.Auth, _batchID)
				if err != nil {
					_logger.Error("Reporting error, err: %v\n", err)
				}
			}
			r.Mux.Unlock()
			break
		}
	}
}

func (r *Reporter) Listening() {
	go r.WatchPartialMatch()
	go r.WatchFullfilMatch()
	go r.WatchReceivedBatchDetails()
	go r.WatchPunishMatcher()
	go r.WatchRemoveBatchOutOfDate()
}

func getWaitingTime(instance *onchain.Onchain) *big.Int {
	waitingTime, err := instance.GetWaitingTime(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		_logger.Error("Get Waiting Time is error, err: %v\n", err)
		return nil
	}
	return waitingTime
}

func (r *Reporter) WatchPartialMatch() {
	logs := make(chan *onchain.OnchainPartialMatch)
	sub, err := r.OnchainInstance.WatchPartialMatch(r.WatchOpts, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		// case vLogs := <-logs:
		case <-logs:
		}
	}
}

func (r *Reporter) WatchFullfilMatch() {
	logs := make(chan *onchain.OnchainFullfilMatch)
	sub, err := r.OnchainInstance.WatchFullfilMatch(r.WatchOpts, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			r.Mux.Lock()
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			r.PendingBatches[id] = time.Now().Unix()
			r.Mux.Unlock()
		}
	}
}

func (r *Reporter) WatchReceivedBatchDetails() {
	logs := make(chan *onchain.OnchainReceivedBatchDetails)
	sub, err := r.OnchainInstance.WatchReceivedBatchDetails(r.WatchOpts, logs)
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
			_logger.Info("Remove batch::%v\n", id.String())
			r.Mux.Lock()
			delete(r.PendingBatches, id)
			r.Mux.Unlock()
		}
	}
}

func (r *Reporter) WatchPunishMatcher() {
	logs := make(chan *onchain.OnchainPunishMatcher)
	sub, err := r.OnchainInstance.WatchPunishMatcher(r.WatchOpts, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		// case vLogs := <-logs:
		case <-logs:
		}
	}
}

func (r *Reporter) WatchRemoveBatchOutOfDate() {
	logs := make(chan *onchain.OnchainRemoveBatchOutOfDate)
	sub, err := r.OnchainInstance.WatchRemoveBatchOutOfDate(r.WatchOpts, logs)
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
			// _logger.Info("Remove batch::%v\n", id.String())
			r.Mux.Lock()
			delete(r.PendingBatches, id)
			r.Mux.Unlock()
		}
	}
}

func (r *Reporter) prepareNonceAndGasPrice(value float64, gasLimit int) {
	nonce, err := r.Client.PendingNonceAt(context.Background(), r.Address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := r.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	r.Auth.Nonce = big.NewInt(int64(nonce))
	r.Auth.GasPrice = gasPrice
	r.Auth.Value = client.EthToWei(big.NewFloat(float64(value)))
	r.Auth.GasLimit = uint64(gasLimit)
}
