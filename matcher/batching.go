package matcher

import (
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/google/uuid"
)

type Batch struct {
	batchID uuid.UUID
	price   int
	amount  int
	side    bool
	orders  []*MatcherOrder
}

func newBatch(_price int, _amount int, _side bool, _orders []*MatcherOrder) *Batch {
	id, _ := uuid.NewRandom()
	return &Batch{
		batchID: id,
		price:   _price,
		amount:  _amount,
		side:    _side,
		orders:  _orders,
	}
}

// I know this function isn't ideal and has too much duplicated code.
// However, the current priority is to get the protocol running.
func (m *Matcher) batching() ([]*Batch, []*Batch) {

	// bid batches
	bidBatches := []*Batch{}
	if len(m.BidOrders) > 0 {
		ord := m.BidOrders[0]
		m.BidOrders = m.BidOrders[1:]
		batch := newBatch(int(ord.Data.Price), int(ord.Data.Amount), ord.Data.Side, []*MatcherOrder{ord})
		cnt := 1
		for cnt <= constants.NO_BATCHES_EACH_TIME && len(m.BidOrders) > 0 {
			ord = m.BidOrders[0]
			if ord.Data.Price == int64(batch.price) {
				m.BidOrders = m.BidOrders[1:]
				batch.amount += int(ord.Data.Amount)
				batch.orders = append(batch.orders, ord)
			} else {
				cnt++
				bidBatches = append(bidBatches, batch)
				batch = newBatch(int(ord.Data.Price), int(ord.Data.Amount), ord.Data.Side, []*MatcherOrder{ord})
			}
		}
		if cnt < constants.NO_BATCHES_EACH_TIME {
			for len(m.BidOrders) > 0 {
				ord = m.BidOrders[0]
				if ord.Data.Price == int64(batch.price) {
					m.BidOrders = m.BidOrders[1:]
					batch.amount += int(ord.Data.Amount)
					batch.orders = append(batch.orders, ord)
				} else {
					break
				}
			}
			bidBatches = append(bidBatches, batch)
		}
	}

	// ask batches
	askBatches := []*Batch{}
	if len(m.AskOrders) > 0 {
		ord := m.AskOrders[0]
		m.AskOrders = m.AskOrders[1:]
		batch := newBatch(int(ord.Data.Price), int(ord.Data.Amount), ord.Data.Side, []*MatcherOrder{ord})
		cnt := 1
		for cnt <= constants.NO_BATCHES_EACH_TIME && len(m.AskOrders) > 0 {
			ord = m.AskOrders[0]
			if ord.Data.Price == int64(batch.price) {
				m.AskOrders = m.AskOrders[1:]
				batch.amount += int(ord.Data.Amount)
				batch.orders = append(batch.orders, ord)
			} else {
				cnt++
				askBatches = append(askBatches, batch)
				batch = newBatch(int(ord.Data.Price), int(ord.Data.Amount), ord.Data.Side, []*MatcherOrder{ord})
			}
		}
		if cnt < constants.NO_BATCHES_EACH_TIME {
			for len(m.AskOrders) > 0 {
				ord = m.AskOrders[0]
				if ord.Data.Price == int64(batch.price) {
					m.AskOrders = m.AskOrders[1:]
					batch.amount += int(ord.Data.Amount)
					batch.orders = append(batch.orders, ord)
				} else {
					break
				}
			}
			askBatches = append(askBatches, batch)
		}
	}

	return bidBatches, askBatches
}
