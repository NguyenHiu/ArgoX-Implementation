package matcher

import (
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/tradeApp"
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

		bidOrder := m.BidOrders[0]
		askOrder := m.AskOrders[0]

		// TODO: Send messages after matching!
		minAmount := bidOrder.Data.Amount
		if minAmount.Cmp(askOrder.Data.Amount) == 1 {
			minAmount = askOrder.Data.Amount
		}

		bidOrder.Data.Amount = new(big.Int).Sub(bidOrder.Data.Amount, minAmount)
		askOrder.Data.Amount = new(big.Int).Sub(askOrder.Data.Amount, minAmount)

		matchPrice := new(big.Int).Div(new(big.Int).Add(bidOrder.Data.Price, askOrder.Data.Price), big.NewInt(2))

		trade := m.NewTrade(bidOrder.Data.From, askOrder.Data.From, matchPrice, minAmount)

		_bidOrder, ok := m.Orders[bidOrder.Data.From]
		if !ok {
			_logger.Debug("can not found bid order\n")
		}

		_askOrder, ok := m.Orders[askOrder.Data.From]
		if !ok {
			_logger.Debug("can not found ask order\n")
		}
		m.ClientConfigs[bidOrder.Owner].TradeChannel.SendNewTrades([]*tradeApp.Trade{trade}, _bidOrder, _askOrder, true)
		m.ClientConfigs[askOrder.Owner].TradeChannel.SendNewTrades([]*tradeApp.Trade{trade}, _bidOrder, _askOrder, false)

		if bidOrder.Data.Amount.Cmp(new(big.Int)) == 0 {
			// m.BidOrders = m.BidOrders[1:]
			for i, _bo := range m.BidOrders {
				if _bo.Data.Equal(bidOrder.Data) {
					m.BidOrders = append(m.BidOrders[:i], m.BidOrders[i+1:]...)
					delete(m.Orders, bidOrder.Data.From)
					break
				}
			}
		}
		if askOrder.Data.Amount.Cmp(new(big.Int)) == 0 {
			// m.AskOrders = m.AskOrders[1:]
			for i, _ao := range m.AskOrders {
				if _ao.Data.Equal(askOrder.Data) {
					m.AskOrders = append(m.AskOrders[:i], m.AskOrders[i+1:]...)
					delete(m.Orders, askOrder.Data.From)
					break
				}
			}
		}
		// <-time.After(time.Millisecond * 500)
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
