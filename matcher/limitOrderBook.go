package matcher

import (
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/app"
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
	if len(m.BidOrders) >= 10 || len(m.AskOrders) >= 10 {
		batches := m.batching()
		for _, batch := range batches {
			m.SendBatch(batch)
		}
	}
}

func (m *Matcher) Log() {
	_logger.Debug("-----------------------------\n")
	_logger.Debug("Local Order Book:\n")
	_logger.Debug("BID:\n")
	for _, order := range m.BidOrders {
		_logger.Debug("\t[%v] %v - %v\n", order.Data.From.String()[:5], order.Data.Price, order.Data.Amount)
	}
	_logger.Debug("---------------\n")
	_logger.Debug("ASK:\n")
	for _, order := range m.AskOrders {
		_logger.Debug("\t[%v] %v - %v\n", order.Data.From.String()[:5], order.Data.Price, order.Data.Amount)
	}
	_logger.Debug("-----------------------------\n")
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
		_logger.Debug("Matching: (%v..., %v), (%v..., %v)\n", m.BidOrders[0].Data.From.String()[:5], m.BidOrders[0].Data.Amount, m.AskOrders[0].Data.From.String()[:5], m.AskOrders[0].Data.Amount)

		// TODO: Send messages after matching!
		minAmount := m.BidOrders[0].Data.Amount
		if minAmount.Cmp(m.AskOrders[0].Data.Amount) == 1 {
			minAmount = m.AskOrders[0].Data.Amount
		}

		m.BidOrders[0].Data.Amount = new(big.Int).Sub(m.BidOrders[0].Data.Amount, minAmount)
		m.AskOrders[0].Data.Amount = new(big.Int).Sub(m.AskOrders[0].Data.Amount, minAmount)

		matchPrice := new(big.Int).Div(new(big.Int).Add(m.BidOrders[0].Data.Price, m.AskOrders[0].Data.Price), big.NewInt(2))

		trade := m.NewTrade(m.BidOrders[0].Data.From, m.AskOrders[0].Data.From, matchPrice, minAmount)

		m.ClientConfigs[m.BidOrders[0].Owner].VerifyChannel.SendNewTrades([]*app.Trade{trade})
		// _logger.Debug("m.AskOrders[0].Owner: %v\n", m.AskOrders[0].Owner)
		m.ClientConfigs[m.AskOrders[0].Owner].VerifyChannel.SendNewTrades([]*app.Trade{trade})

		if m.BidOrders[0].Data.Amount.Cmp(new(big.Int)) == 0 {
			m.BidOrders = m.BidOrders[1:]
		}
		if m.AskOrders[0].Data.Amount.Cmp(new(big.Int)) == 0 {
			m.AskOrders = m.AskOrders[1:]
		}
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
