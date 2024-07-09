package worker

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/orderClient"
	"github.com/google/uuid"
)

func (w *Worker) SubmitMatchEvent(bidBatchID, askBatchID uuid.UUID) {
	w.prepareNonceAndGasPrice(0, 300000)

	b, err := bidBatchID.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	var bArray [16]byte
	copy(bArray[:], b)
	a, err := askBatchID.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	var aArray [16]byte
	copy(aArray[:], a)

	_logger.Info("Submit a match event\n")
	_, err = w.OnchainIstance.Matching(w.Auth, bArray, aArray)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Worker) Listening() {
	go w.WatchAcceptBatch()
}

func (w *Worker) WatchAcceptBatch() {
	logs := make(chan *onchain.OnchainAcceptBatch)
	sub, err := w.OnchainIstance.WatchAcceptBatch(w.WatchOpts, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			batchID, err := uuid.FromBytes(vLogs.Arg0[:])
			if err != nil {
				log.Fatal(err)
			}

			_logger.Info("Receive a batch::%v\n", batchID)

			w.addBatch(&Batch{
				BatchID: batchID,
				Price:   int(vLogs.Arg1.Int64()),
				Amount:  int(vLogs.Arg2.Int64()),
				Side:    vLogs.Arg3,
			})
		}
	}
}

func (w *Worker) WatchFullfilMatch() {
	logs := make(chan *onchain.OnchainFullfilMatch)
	sub, err := w.OnchainIstance.WatchFullfilMatch(w.WatchOpts, logs)
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
			// _logger.Info("[Fullfill] Batch::%v\n", id.String())
			batch, ok := w.Batches[id]
			if ok {
				w.Mux.Lock()
				if batch.Side == constants.BID {
					for i, b := range w.BidBatches {
						if b.Equal(batch) {
							w.BidBatches = append(w.BidBatches[:i], w.BidBatches[i+1:]...)
							break
						}
					}
				} else {
					for i, b := range w.AskBatches {
						if b.Equal(batch) {
							w.AskBatches = append(w.AskBatches[:i], w.AskBatches[i+1:]...)
							break
						}
					}
				}
				w.Mux.Unlock()
			}
		}
	}
}

func (w *Worker) prepareNonceAndGasPrice(value float64, gasLimit int) {
	nonce, err := w.Client.PendingNonceAt(context.Background(), w.Address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := w.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	w.Auth.Nonce = big.NewInt(int64(nonce))
	w.Auth.GasPrice = gasPrice
	w.Auth.Value = orderClient.EthToWei(big.NewFloat(float64(value)))
	w.Auth.GasLimit = uint64(gasLimit)
}
