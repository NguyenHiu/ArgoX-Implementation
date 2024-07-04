package matcher

import (
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
)

func (m *Matcher) SendBatch(batch *Batch) {
	m.Mux.Lock()
	m.Batches[batch.BatchID] = batch
	m.Mux.Unlock()

	// Send instantly
	orders := []*supermatcher.ExpandOrder{}
	for _, order := range batch.Orders {
		orders = append(orders, &supermatcher.ExpandOrder{
			ShadowOrder:   (*supermatcher.ShadowOrder)(order.ShadowOrder),
			Trades:        order.Trades,
			OriginalOrder: order.OriginalOrder,
		})
	}
	m.SuperMatcherInstance.SendBatch(&supermatcher.Batch{
		BatchID:   batch.BatchID,
		Price:     batch.Price,
		Amount:    batch.Amount,
		Side:      batch.Side,
		Orders:    orders,
		Owner:     batch.Owner,
		Signature: batch.Signature,
	})

	m.SuperMatcherInstance.Process()
}
