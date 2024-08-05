package matcher

import (
	"fmt"
	"log"
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

		// bid check
		totalTrade := new(big.Int)
		for _, t := range m.mappingBidtoTrade[bidOrder.Data.From] {
			totalTrade.Add(totalTrade, t.Amount)
		}
		if new(big.Int).Add(bidOrder.Data.Amount, totalTrade).Cmp(m.Orders[bidOrder.Data.From].Amount) == 1 {
			log.Fatalf("bid order has invlaid trade slice\ntotalTrade: %v\nbidOrder.Data.Amount: %v\nm.Orders[bidOrder.Data.From].Amount: %v\n", totalTrade, bidOrder.Data.Amount, m.Orders[bidOrder.Data.From].Amount)
		}
		// ask check
		totalTrade = new(big.Int)
		for _, t := range m.mappingBidtoTrade[askOrder.Data.From] {
			totalTrade.Add(totalTrade, t.Amount)
		}
		if new(big.Int).Add(askOrder.Data.Amount, totalTrade).Cmp(m.Orders[askOrder.Data.From].Amount) == 1 {
			log.Fatalf("bid order has invlaid trade slice\ntotalTrade: %v\naskOrder.Data.Amount: %v\nm.Orders[askOrder.Data.From].Amount: %v\n", totalTrade, askOrder.Data.Amount, m.Orders[askOrder.Data.From].Amount)
		}

		// Get minimize amount among bid & ask order
		minAmount := new(big.Int).Set(bidOrder.Data.Amount)
		if minAmount.Cmp(askOrder.Data.Amount) == 1 {
			minAmount = new(big.Int).Set(askOrder.Data.Amount)
		}

		// status:
		//   - 0: not changed, but failed
		//   - 1: changed, but failed
		// 	 - 2: changed, and success
		_bidChange, _askChange, _bidLeftAmount, _askLeftAmount := m.SuperMatcherInstance.MatchAnOrder(
			bidOrder.Data.From, m.Orders[bidOrder.Data.From].Amount,
			askOrder.Data.From, m.Orders[askOrder.Data.From].Amount,
			minAmount,
		)

		_logger.Debug("bid's id: %v\n", bidOrder.Data.From)
		_logger.Debug("bid's amount: %v\n", bidOrder.Data.Amount)
		_logger.Debug("ask's id: %v\n", askOrder.Data.From)
		_logger.Debug("ask's amount: %v\n", askOrder.Data.Amount)
		_logger.Debug("minAmount: %v\n", minAmount)

		if _bidChange != 2 || _askChange != 2 {
			if _bidChange == 1 {
				if _bidLeftAmount.Cmp(new(big.Int)) == 0 {
					_logger.Debug("ID: %v - Order was cleared\n", bidOrder.Data.From)
					m.BidOrders = m.BidOrders[1:]
					// delete(m.Orders, bidOrder.Data.From)
				} else {
					bidOrder.Data.Amount = new(big.Int).Set(_bidLeftAmount)
					_logger.Debug("Match fail, ID: %v, , bid left amount: %v\n", bidOrder.Data.From, _bidLeftAmount)
				}
			}

			if _askChange == 1 {
				if _askLeftAmount.Cmp(new(big.Int)) == 0 {
					_logger.Debug("ID: %v - Order was cleared\n", askOrder.Data.From)
					m.AskOrders = m.AskOrders[1:]
					// delete(m.Orders, askOrder.Data.From)
				} else {
					askOrder.Data.Amount = new(big.Int).Set(_askLeftAmount)
					_logger.Debug("Match fail, ID: %v, ask left amount: %v\n", askOrder.Data.From, _askLeftAmount)
				}
			}

			continue
		}

		if _bidLeftAmount.Cmp(new(big.Int).Sub(bidOrder.Data.Amount, minAmount)) == 1 {
			s := "bid\n"
			for _, o := range m.mappingBidtoTrade[bidOrder.Data.From] {
				s = fmt.Sprintf("%vexecuted trade: %v", s, o.Amount)
			}

			log.Fatalf("INVALID MATCHING BID\nID: %v\n_bidLeftAmount: %v\nbidOrder.Data.Amount: %v\nminAmount: %v\ntraded: %v\noriginal_amount: %v\n", bidOrder.Data.From, _bidLeftAmount, bidOrder.Data.Amount, minAmount, s, m.Orders[bidOrder.Data.From].Amount)
		}
		if _askLeftAmount.Cmp(new(big.Int).Sub(askOrder.Data.Amount, minAmount)) == 1 {
			s := "ask\n"
			for _, o := range m.mappingBidtoTrade[bidOrder.Data.From] {
				s = fmt.Sprintf("%vexecuted trade: %v", s, o.Amount)
			}
			log.Fatalf("INVALID MATCHING ASK\nID: %v\n_askLeftAmount: %v\naskOrder.Data.Amount: %v\nminAmount: %v\ntraded: %v\noriginal_amount: %v\n", askOrder.Data.From, _askLeftAmount, askOrder.Data.Amount, minAmount, s, m.Orders[askOrder.Data.From].Amount)
		}

		bidOrder.Data.Amount = new(big.Int).Set(_bidLeftAmount)
		askOrder.Data.Amount = new(big.Int).Set(_askLeftAmount)
		m.TotalMatchedAmountLocal.Add(m.TotalMatchedAmountLocal, minAmount)
		m.TotalMatchedAmountLocal.Add(m.TotalMatchedAmountLocal, minAmount)
		m.TotalTimeLocal += time.Now().Unix() - m.CreateTime[bidOrder.Data.From]
		m.TotalTimeLocal += time.Now().Unix() - m.CreateTime[askOrder.Data.From]

		matchPrice := new(big.Int).Div(new(big.Int).Add(bidOrder.Data.Price, askOrder.Data.Price), big.NewInt(2))
		// m.PriceCurveLocal = append(m.PriceCurveLocal, matchPrice)
		m.TotalProfitLocal.Add(m.TotalProfitLocal, new(big.Int).Mul(new(big.Int).Sub(bidOrder.Data.Price, askOrder.Data.Price), minAmount))
		m.TotalRawProfitLocal.Add(m.TotalRawProfitLocal, new(big.Int).Sub(bidOrder.Data.Price, askOrder.Data.Price))
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
			m.NumberOfMatchedOrder += 1
			for i, _bo := range m.BidOrders {
				if _bo.Data.Equal(bidOrder.Data) {
					m.BidOrders = append(m.BidOrders[:i], m.BidOrders[i+1:]...)
					// delete(m.Orders, bidOrder.Data.From)
					break
				}
			}
			m.OrderStatusMapping[bidOrder.Data.From] = FULFILED
		} else {
			m.OrderStatusMapping[bidOrder.Data.From] = fmt.Sprintf("PARTIAL_MATCH,%v", m.SuperMatcherInstance.GetRemainingAmount(bidOrder.Data.From))
		}
		if askOrder.Data.Amount.Cmp(new(big.Int)) == 0 {
			m.NumberOfMatchedOrder += 1
			for i, _ao := range m.AskOrders {
				if _ao.Data.Equal(askOrder.Data) {
					m.AskOrders = append(m.AskOrders[:i], m.AskOrders[i+1:]...)
					// delete(m.Orders, askOrder.Data.From)
					break
				}
			}
			m.OrderStatusMapping[askOrder.Data.From] = FULFILED
		} else {
			m.OrderStatusMapping[askOrder.Data.From] = fmt.Sprintf("PARTIAL_MATCH,%v", m.SuperMatcherInstance.GetRemainingAmount(askOrder.Data.From))
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
