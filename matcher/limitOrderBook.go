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
	m.matching()
	// if len(m.BidOrders) >= 30 || len(m.AskOrders) >= 30 {
	// 	batches := m.batching()
	// 	for _, batch := range batches {
	// 		m.SendBatch(batch)
	// 	}
	// }
}

func (m *Matcher) Log() {
	//IMHERETODEBUG_logger.Debug("-----------------------------\n")
	//IMHERETODEBUG_logger.Debug("Local Order Book:\n")
	//IMHERETODEBUG_logger.Debug("BID:\n")
	for _, order := range m.BidOrders {
		//IMHERETODEBUG_logger.Debug("\t[%v] %v - %v\n", order.Data.From.String()[:5], order.Data.Price, order.Data.Amount)
	}
	//IMHERETODEBUG_logger.Debug("---------------\n")
	//IMHERETODEBUG_logger.Debug("ASK:\n")
	for _, order := range m.AskOrders {
		//IMHERETODEBUG_logger.Debug("\t[%v] %v - %v\n", order.Data.From.String()[:5], order.Data.Price, order.Data.Amount)
	}
	//IMHERETODEBUG_logger.Debug("-----------------------------\n")
}

func (m *Matcher) matching() {
	m.Mux.Lock()
	defer m.Mux.Unlock()

	// naive matching
	for m.canMatch() {
		//IMHERETODEBUG_logger.Debug("Matching: (%v..., %v), (%v..., %v)\n", m.BidOrders[0].Data.From.String()[:5], m.BidOrders[0].Data.Amount, m.AskOrders[0].Data.From.String()[:5], m.AskOrders[0].Data.Amount)

		bidOrder := m.BidOrders[0]
		askOrder := m.AskOrders[0]

		// Get minimize amount among bid & ask order
		minAmount := new(big.Int).Set(bidOrder.Data.Amount)
		if minAmount.Cmp(m.AskOrders[0].Data.Amount) == 1 {
			minAmount = new(big.Int).Set(askOrder.Data.Amount)
		}

		_isValid, _bidLeftAmount, _askLeftAmount := m.SuperMatcherInstance.MatchAnOrder(
			bidOrder.Data.From, minAmount, m.Orders[bidOrder.Data.From].Amount,
			askOrder.Data.From, minAmount, m.Orders[askOrder.Data.From].Amount,
		)

		if !_isValid {
			if _bidLeftAmount.Cmp(new(big.Int)) == 0 {
				_bid := bidOrder
				m.BidOrders = m.BidOrders[1:]
				delete(m.Orders, _bid.Data.From)
			} else {
				bidOrder.Data.Amount = new(big.Int).Set(_bidLeftAmount)
			}
			if _askLeftAmount.Cmp(new(big.Int)) == 0 {
				_ask := askOrder
				m.AskOrders = m.AskOrders[1:]
				delete(m.Orders, _ask.Data.From)
			} else {
				askOrder.Data.Amount = new(big.Int).Set(_askLeftAmount)
			}
			continue
		}

		if _bidLeftAmount.Cmp(new(big.Int).Sub(bidOrder.Data.Amount, minAmount)) != 0 {
			log.Fatal("INVALID MATCHING BID")
		}
		if _askLeftAmount.Cmp(new(big.Int).Sub(askOrder.Data.Amount, minAmount)) != 0 {
			log.Fatal("INVALID MATCHING ASK")
		}

		bidOrder.Data.Amount = new(big.Int).Set(_bidLeftAmount)
		askOrder.Data.Amount = new(big.Int).Set(_askLeftAmount)

		// //IMHERETODEBUG_logger.Debug("Matched, amount: %v\n", minAmount)
		// //IMHERETODEBUG_logger.Debug("Matched, amount: %v\n", minAmount)
		// //IMHERETODEBUG_logger.Debug("Time: %v\n", time.Now().Unix()-m.CreateTime[bidOrder.Data.From])
		// //IMHERETODEBUG_logger.Debug("Time: %v\n", time.Now().Unix()-m.CreateTime[askOrder.Data.From])
		m.TotalMatchedAmountLocal.Add(m.TotalMatchedAmountLocal, minAmount)
		m.TotalMatchedAmountLocal.Add(m.TotalMatchedAmountLocal, minAmount)
		m.TotalTimeLocal += time.Now().Unix() - m.CreateTime[bidOrder.Data.From]
		m.TotalTimeLocal += time.Now().Unix() - m.CreateTime[askOrder.Data.From]

		fmt.Printf("{\"ID\": \"%v\", \"Amount\": %v},", bidOrder.Data.From, minAmount)
		fmt.Printf("{\"ID\": \"%v\", \"Amount\": %v},", askOrder.Data.From, minAmount)

		// IMHERETODEBUG_logger.Debug("[DEBUG FLAG] %v - %v\n", bidOrder.Data.From, minAmount)
		// IMHERETODEBUG_logger.Debug("[DEBUG FLAG] %v - %v\n", askOrder.Data.From, minAmount)

		matchPrice := new(big.Int).Div(new(big.Int).Add(bidOrder.Data.Price, askOrder.Data.Price), big.NewInt(2))
		// m.PriceCurveLocal = append(m.PriceCurveLocal, matchPrice)
		m.TotalProfitLocal.Add(m.TotalProfitLocal, new(big.Int).Mul(new(big.Int).Sub(bidOrder.Data.Price, askOrder.Data.Price), minAmount))
		m.CurrentPrice = new(big.Int).Set(matchPrice)

		trade := m.NewTrade(bidOrder.Data.From, askOrder.Data.From, matchPrice, minAmount)

		_bidOrder, ok := m.Orders[bidOrder.Data.From]
		if !ok {
			//IMHERETODEBUG_logger.Debug("can not found bid order\n")
		}

		_askOrder, ok := m.Orders[askOrder.Data.From]
		if !ok {
			//IMHERETODEBUG_logger.Debug("can not found ask order\n")
		}
		m.ClientConfigs[bidOrder.Owner].TradeChannel.SendNewTrades([]*tradeApp.Trade{trade}, _bidOrder, _askOrder, true)
		m.ClientConfigs[askOrder.Owner].TradeChannel.SendNewTrades([]*tradeApp.Trade{trade}, _bidOrder, _askOrder, false)

		if bidOrder.Data.Amount.Cmp(new(big.Int)) == 0 {
			m.NumberOfMatchedOrder += 1
			for i, _bo := range m.BidOrders {
				if _bo.Data.Equal(bidOrder.Data) {
					m.BidOrders = append(m.BidOrders[:i], m.BidOrders[i+1:]...)
					delete(m.Orders, bidOrder.Data.From)
					break
				}
			}
		}
		if askOrder.Data.Amount.Cmp(new(big.Int)) == 0 {
			m.NumberOfMatchedOrder += 1
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
