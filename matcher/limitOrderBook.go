package matcher

import (
	"fmt"
	"log"

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

	if len(m.BidOrders) == 4 {
		fmt.Println("Order Book:")
		for _, order := range m.BidOrders {
			fmt.Printf("%v, %v\n", order.Data.Price, order.Data.Amount)
		}
		fmt.Println()
	}
}

func (m *Matcher) matching() bool {
	if !m.canMatch() {
		return false
	}

	// naive matching
	for m.canMatch() {
		log.Printf("Matching (%v..., %v..., %v)\n", m.BidOrders[0].Data.OrderID[:5], m.AskOrders[0].Data.OrderID[:5], m.BidOrders[0].Data.Amount)
		m.ClientConfigs[m.BidOrders[0].Owner].VerifyChannel.UpdateExistedOrder(
			m.BidOrders[0].Data.OrderID, app.OrderUpdatedInfo{
				Status:        "M",
				MatchedAmount: m.BidOrders[0].Data.Amount,
			},
		)
		m.ClientConfigs[m.AskOrders[0].Owner].VerifyChannel.UpdateExistedOrder(
			m.AskOrders[0].Data.OrderID, app.OrderUpdatedInfo{
				Status:        "M",
				MatchedAmount: m.AskOrders[0].Data.Amount,
			},
		)
		m.BidOrders = m.BidOrders[1:]
		m.AskOrders = m.AskOrders[1:]
	}
	return true
}

func (m *Matcher) canMatch() bool {
	if len(m.BidOrders) == 0 || len(m.AskOrders) == 0 {
		return false
	}
	return m.BidOrders[0].Data.Price >= m.AskOrders[0].Data.Price
}

func addAccordingTheOrder(order *MatcherOrder, orders []*MatcherOrder) []*MatcherOrder {
	l := len(orders)
	if l == 0 {
		orders = append(orders, order)
	} else if l == 1 {
		if (order.Data.Side == constants.BID && order.Data.Price > orders[0].Data.Price) ||
			(order.Data.Side == constants.ASK && order.Data.Price < orders[0].Data.Price) {
			orders = append([]*MatcherOrder{order}, orders...)
		} else {
			orders = append(orders, order)
		}
	} else {
		for i := 0; i < l; i++ {
			if (order.Data.Side == constants.BID && order.Data.Price > orders[i].Data.Price) ||
				(order.Data.Side == constants.ASK && order.Data.Price < orders[i].Data.Price) {
				fmt.Println("Ahhh, you touched me?")
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
