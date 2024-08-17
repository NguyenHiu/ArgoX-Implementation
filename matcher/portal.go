package matcher

import (
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/google/uuid"
)

func (m *Matcher) SendBatch(batch *Batch) {
	_logger.Debug("[%v] Send batch::%v, amount: %v, side: %v, price: %v\n", m.Address.String()[:5], batch.BatchID.String()[:5], batch.Amount, batch.Side, batch.Price)

	// Send instantly
	orders := []*supermatcher.ExpandOrder{}
	for _, order := range batch.Orders {
		orders = append(orders, &supermatcher.ExpandOrder{
			ShadowOrder: &supermatcher.ShadowOrder{
				Price:  order.ShadowOrder.Price,
				Amount: new(big.Int).Set(order.ShadowOrder.Amount),
				Side:   order.ShadowOrder.Side,
				From:   order.OriginalOrder.OrderID,
			},
			Trades:        order.Trades,
			OriginalOrder: order.OriginalOrder,
		})
	}
	_status, _validOrders := m.SuperMatcherInstance.AddBatch(&supermatcher.Batch{
		BatchID:   batch.BatchID,
		Price:     batch.Price,
		Amount:    new(big.Int).Set(batch.Amount),
		Side:      batch.Side,
		Orders:    orders,
		Owner:     batch.Owner,
		Signature: batch.Signature,
	})

	retry := true
	retryTime := 0
	for retry {
		if retryTime > 10 {
			_logger.Error("SendBatch retry more than 10 times\n")
			return
		}

		if _status == "OK" {
			m.Mux.Lock()
			m.Batches[batch.BatchID] = batch
			m.Mux.Unlock()
			m.SuperMatcherInstance.Process()
			retry = false
			_logger.Debug("Batch::%v (OK)\n", batch.BatchID.String())
		} else if _status == "RESIGN" && len(_validOrders) != 0 {
			retry = true
			retryTime += 1

			_amount := new(big.Int)
			for _, _order := range _validOrders {
				_amount.Add(_amount, _order.ShadowOrder.Amount)
			}

			_orders := []*ExpandOrder{}
			for _, _order := range _validOrders {
				_orders = append(_orders, &ExpandOrder{
					ShadowOrder: &ShadowOrder{
						Price:  _order.ShadowOrder.Price,
						Amount: new(big.Int).Set(_order.ShadowOrder.Amount),
						Side:   _order.ShadowOrder.Side,
						From:   _order.ShadowOrder.From,
					},
					Trades:        _order.Trades,
					OriginalOrder: _order.OriginalOrder,
				})
			}
			_id, _ := uuid.NewRandom()
			_newBatch := &Batch{
				BatchID: _id,
				Price:   _validOrders[0].OriginalOrder.Price,
				Amount:  new(big.Int).Set(_amount),
				Side:    _validOrders[0].OriginalOrder.Side,
				Orders:  _orders,
				Owner:   m.Address,
			}
			_newBatch.Sign(m.PrivateKey)
			_logger.Debug("Batch::%v (RESIGN) --> Batch::%v\n", batch.BatchID.String(), _id.String())
			_logger.Debug(
				"[%v] Send batch::%v, amount: %v, side: %v, price: %v\n",
				m.Address.String()[:5],
				_newBatch.BatchID.String()[:5],
				_newBatch.Amount,
				_newBatch.Side,
				_newBatch.Price,
			)
			_status, _validOrders = m.SuperMatcherInstance.AddBatch(&supermatcher.Batch{
				BatchID:   _newBatch.BatchID,
				Price:     _newBatch.Price,
				Amount:    new(big.Int).Set(_newBatch.Amount),
				Side:      _newBatch.Side,
				Orders:    _validOrders,
				Owner:     _newBatch.Owner,
				Signature: _newBatch.Signature,
			})
		} else {
			_logger.Debug("Batch::%v is empty (cook)\n", batch.BatchID.String())
			retry = false
		}
	}

}
