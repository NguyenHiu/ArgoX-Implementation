package worker

import (
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/google/uuid"
)

type Batch struct {
	BatchID uuid.UUID
	Price   int
	Amount  int
	Side    bool
}

func (b *Batch) Equal(_b *Batch) bool {
	return b.BatchID == _b.BatchID &&
		b.Price == _b.Price &&
		b.Amount == _b.Amount &&
		b.Side == _b.Side
}

func (w *Worker) addBatch(newBatch *Batch) {
	w.Mux.Lock()
	if newBatch.Side == constants.ASK {
		w.AskBatches = addBatchAccording(newBatch, w.AskBatches)
	} else {
		w.BidBatches = addBatchAccording(newBatch, w.BidBatches)
	}
	w.Mux.Unlock()

	w.Batches[newBatch.BatchID] = newBatch

	w.matching()
}

func addBatchAccording(newBatch *Batch, batches []*Batch) []*Batch {
	l := len(batches)
	if l == 0 {
		batches = append(batches, newBatch)
	} else if l == 1 {
		if (newBatch.Side == constants.ASK && newBatch.Price < batches[0].Price) ||
			(newBatch.Side == constants.BID && newBatch.Price > batches[0].Price) {
			batches = append([]*Batch{newBatch}, batches...)
		} else {
			batches = append(batches, newBatch)
		}
	} else {
		for i := 0; i < l; i++ {
			if (newBatch.Side == constants.ASK && newBatch.Price < batches[i].Price) ||
				(newBatch.Side == constants.BID && newBatch.Price > batches[i].Price) {
				batches = append(batches, nil)
				copy(batches[i+1:], batches[i:])
				batches[i] = newBatch
				return batches
			}
		}
		batches = append(batches, newBatch)
	}
	return batches
}

func (w *Worker) matching() {
	if len(w.AskBatches) == 0 || len(w.BidBatches) == 0 {
		return
	}

	// w.log()

	w.Mux.Lock()
	defer w.Mux.Unlock()

	i, j := 0, 0
	for ; w.canMatch(i, j); i++ {
		for ; w.canMatch(i, j); j++ {
			if w.BidBatches[i].Amount == w.AskBatches[j].Amount {
				w.SubmitMatchEvent(w.BidBatches[i].BatchID, w.AskBatches[j].BatchID)
				w.BidBatches = append(w.BidBatches[:i], w.BidBatches[i+1:]...)
				w.AskBatches = append(w.AskBatches[:j], w.AskBatches[j+1:]...)
			}
		}
		j = 0
	}
}

func (w *Worker) canMatch(i, j int) bool {
	if i >= len(w.BidBatches) || j >= len(w.AskBatches) {
		return false
	}
	return w.BidBatches[i].Price >= w.AskBatches[j].Price
}

func (w *Worker) log() {
	//IMHERETODEBUG_logger.Debug("======ASK===================================\n")
	for _, batch := range w.AskBatches {
		//IMHERETODEBUG_logger.Debug("\t%v - %v\n", batch.Price, batch.Amount)
	}
	//IMHERETODEBUG_logger.Debug("======BID===================================\n")
	for _, batch := range w.BidBatches {
		//IMHERETODEBUG_logger.Debug("\t%v - %v\n", batch.Price, batch.Amount)
	}
	//IMHERETODEBUG_logger.Debug("\n")
}
