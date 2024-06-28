package matcher

import (
	"github.com/NguyenHiu/lightning-exchange/constants"
)

// implement limit order book logic

// bid: 1 --> 2
// ask: 2 --> 1
func (m *Matcher) addOrder(order *MatcherOrder) {
	if order.Data.Side == constants.BID {
		m.BidOrders = addAccordingTheOrder(order, m.BidOrders)
	} else {
		m.AskOrders = addAccordingTheOrder(order, m.AskOrders)
	}
	m.matching()
}

func (m *Matcher) matching() bool {
	if !m.canMatch() {
		return false
	}

	m.Mux.Lock()
	defer m.Mux.Unlock()

	// <-time.After(time.Second)

	// naive matching
	for m.canMatch() {
		_logger.Debug("Matching (%v..., %v..., %v)\n", m.BidOrders[0].Data.OrderID.String()[:5], m.AskOrders[0].Data.OrderID.String()[:5], m.BidOrders[0].Data.Amount)

		m.SendNewMessage(m.BidOrders[0].Owner, m.BidOrders[0].Data.OrderID, m.BidOrders[0].Data.Amount, 'M')
		m.SendNewMessage(m.AskOrders[0].Owner, m.AskOrders[0].Data.OrderID, m.AskOrders[0].Data.Amount, 'M')

		m.BidOrders = m.BidOrders[1:]
		m.AskOrders = m.AskOrders[1:]
	}
	return true
}

func (m *Matcher) canMatch() bool {
	if len(m.BidOrders) == 0 || len(m.AskOrders) == 0 {
		return false
	}
	return m.BidOrders[0].Data.Price.Cmp(m.AskOrders[0].Data.Price) != -1
}

func addAccordingTheOrder(order *MatcherOrder, orders []*MatcherOrder) []*MatcherOrder {
	l := len(orders)
	if l == 0 {
		orders = append(orders, order)
	} else if l == 1 {
		if (order.Data.Side == constants.BID && order.Data.Price.Cmp(orders[0].Data.Price) == 1) ||
			(order.Data.Side == constants.ASK && order.Data.Price.Cmp(orders[0].Data.Price) == -1) {
			orders = append([]*MatcherOrder{order}, orders...)
		} else {
			orders = append(orders, order)
		}
	} else {
		for i := 0; i < l; i++ {
			if (order.Data.Side == constants.BID && order.Data.Price.Cmp(orders[i].Data.Price) == 1) ||
				(order.Data.Side == constants.ASK && order.Data.Price.Cmp(orders[i].Data.Price) == -1) {
				orders = append(orders, nil)
				copy(orders[i+1:], orders[i:])
				orders[i] = order
				return orders
			}
		}
		orders = append(orders, order)
	}

	return orders
}
