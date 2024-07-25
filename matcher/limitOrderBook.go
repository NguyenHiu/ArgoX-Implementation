package matcher

import (
	"math/big"
	"time"

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
	if len(m.BidOrders) >= 7 || len(m.AskOrders) >= 7 {
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

func (m *Matcher) matching() {
	m.Mux.Lock()
	defer m.Mux.Unlock()

	// naive matching
	for m.canMatch() {
		_logger.Debug("Matching: (%v..., %v), (%v..., %v)\n", m.BidOrders[0].Data.From.String()[:5], m.BidOrders[0].Data.Amount, m.AskOrders[0].Data.From.String()[:5], m.AskOrders[0].Data.Amount)

		bidOrder := m.BidOrders[0]
		askOrder := m.AskOrders[0]

		// Get minimize amount among bid & ask order
		minAmount := new(big.Int).Set(bidOrder.Data.Amount)
		if minAmount.Cmp(askOrder.Data.Amount) == 1 {
			minAmount = new(big.Int).Set(askOrder.Data.Amount)
		}

		// Check if bidOrder is valid and didnt be matched
		if leftAmount := m.SuperMatcherInstance.GetLeftAmount(bidOrder.Data.From); leftAmount.Cmp(big.NewInt(-1)) != 0 &&
			leftAmount.Cmp(bidOrder.Data.Amount) == -1 {
			m.BidOrders = m.BidOrders[1:]
			continue
		}
		// Check if askOrder is valid and didnt be matched
		if leftAmount := m.SuperMatcherInstance.GetLeftAmount(askOrder.Data.From); leftAmount.Cmp(big.NewInt(-1)) != 0 &&
			leftAmount.Cmp(askOrder.Data.Amount) == -1 {
			m.AskOrders = m.AskOrders[1:]
			continue
		}

		_logger.Debug("Matched, amount: %v\n", minAmount)
		_logger.Debug("Matched, amount: %v\n", minAmount)
		_logger.Debug("Time: %v\n", time.Now().Unix()-m.CreateTime[bidOrder.Data.From])
		_logger.Debug("Time: %v\n", time.Now().Unix()-m.CreateTime[askOrder.Data.From])
		m.TotalMatchedAmountLocal.Add(m.TotalMatchedAmountLocal, minAmount)
		m.TotalMatchedAmountLocal.Add(m.TotalMatchedAmountLocal, minAmount)
		m.TotalTimeLocal += time.Now().Unix() - m.CreateTime[bidOrder.Data.From]
		m.TotalTimeLocal += time.Now().Unix() - m.CreateTime[askOrder.Data.From]

		bidOrder.Data.Amount.Sub(bidOrder.Data.Amount, minAmount)
		askOrder.Data.Amount.Sub(askOrder.Data.Amount, minAmount)

		if !m.SuperMatcherInstance.MatchAnOrder(bidOrder.Data.From, bidOrder.Data.Amount) {
			_logger.Error("invalid action: matching an invalid order (bid)\n")
		}
		if !m.SuperMatcherInstance.MatchAnOrder(askOrder.Data.From, askOrder.Data.Amount) {
			_logger.Error("invalid action: matching an invalid order (ask)\n")
		}

		matchPrice := new(big.Int).Div(new(big.Int).Add(bidOrder.Data.Price, askOrder.Data.Price), big.NewInt(2))
		// m.PriceCurveLocal = append(m.PriceCurveLocal, matchPrice)
		m.CurrentPrice = new(big.Int).Set(matchPrice)

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
			for i, _bo := range m.BidOrders {
				if _bo.Data.Equal(bidOrder.Data) {
					m.BidOrders = append(m.BidOrders[:i], m.BidOrders[i+1:]...)
					delete(m.Orders, bidOrder.Data.From)
					break
				}
			}
		}
		if askOrder.Data.Amount.Cmp(new(big.Int)) == 0 {
			for i, _ao := range m.AskOrders {
				if _ao.Data.Equal(askOrder.Data) {
					m.AskOrders = append(m.AskOrders[:i], m.AskOrders[i+1:]...)
					delete(m.Orders, askOrder.Data.From)
					break
				}
			}
		}
	}
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
