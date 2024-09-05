package tradeApp

import (
	"encoding/binary"
	"io"

	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

type TradeAppData struct {
	Orders        []*Order
	Trades        []*Trade
	OrdersMapping map[uuid.UUID]*Order
	TradesMapping map[uuid.UUID]*Trade
	BidToTrade    map[uuid.UUID][]*Trade
	AskToTrade    map[uuid.UUID][]*Trade
}

/**
 * Encode encodes app data ([]byte) onto an io.Writer.
 * Format: <NO Orders> <Each Order>... <No message lists> [<No messages in each list> <Each message>]...
 */
func (d *TradeAppData) Encode(w io.Writer) error {
	// No Orders
	if err := binary.Write(w, binary.BigEndian, uint8(len(d.Orders))); err != nil {
		return err
	}

	// Each Order
	for _, order := range d.Orders {
		if err := binary.Write(w, binary.BigEndian, order.Encode_TransferLightning()); err != nil {
			return err
		}
	}

	// No Trades
	if err := binary.Write(w, binary.BigEndian, uint8(len(d.Trades))); err != nil {
		return err
	}

	// Each Trades
	for _, trade := range d.Trades {
		data, err := trade.Encode_TransferLightning()
		if err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, data); err != nil {
			return err
		}
	}

	return nil
}

func (d *TradeAppData) Clone() channel.Data {
	cloned := &TradeAppData{
		Orders:        make([]*Order, len(d.Orders)),
		Trades:        make([]*Trade, len(d.Trades)),
		OrdersMapping: make(map[uuid.UUID]*Order),
		TradesMapping: make(map[uuid.UUID]*Trade),
		BidToTrade:    make(map[uuid.UUID][]*Trade),
		AskToTrade:    make(map[uuid.UUID][]*Trade),
	}

	// Clone Orders
	for i, order := range d.Orders {
		cloned.Orders[i] = order.Clone() // Assuming Order has a Clone method
	}

	// Clone Trades
	for i, trade := range d.Trades {
		cloned.Trades[i] = trade.Clone() // Assuming Trade has a Clone method
	}

	// Clone OrdersMapping
	for key, order := range d.OrdersMapping {
		cloned.OrdersMapping[key] = order.Clone() // Assuming Order has a Clone method
	}

	// Clone TradesMapping
	for key, trade := range d.TradesMapping {
		cloned.TradesMapping[key] = trade.Clone() // Assuming Trade has a Clone method
	}

	// Clone BidToTrade
	for key, trades := range d.BidToTrade {
		clonedTrades := make([]*Trade, len(trades))
		for i, trade := range trades {
			clonedTrades[i] = trade.Clone() // Assuming Trade has a Clone method
		}
		cloned.BidToTrade[key] = clonedTrades
	}

	// Clone AskToTrade
	for key, trades := range d.AskToTrade {
		clonedTrades := make([]*Trade, len(trades))
		for i, trade := range trades {
			clonedTrades[i] = trade.Clone() // Assuming Trade has a Clone method
		}
		cloned.AskToTrade[key] = clonedTrades
	}

	return cloned
}

func (d *TradeAppData) SendNewTrades(trades []*Trade, bidOrder, askOrder *Order, isBid bool) {
	for _, trade := range trades {
		if trade.TradeID == EndID {
			d.Trades = append(d.Trades, trade)
			d.TradesMapping[trade.TradeID] = trade
			continue
		}

		_, ok := d.TradesMapping[trade.TradeID]
		if !ok {
			d.Trades = append(d.Trades, trade)
			d.TradesMapping[trade.TradeID] = trade
			d.AskToTrade[trade.AskOrder] = append(d.AskToTrade[trade.AskOrder], trade)
			d.BidToTrade[trade.BidOrder] = append(d.BidToTrade[trade.BidOrder], trade)
		}

		if isBid {
			if _, ok := d.OrdersMapping[trade.BidOrder]; !ok {
				d.Orders = append(d.Orders, bidOrder)
				d.OrdersMapping[trade.BidOrder] = bidOrder
			}
		} else {
			if _, ok := d.OrdersMapping[trade.AskOrder]; !ok {
				d.Orders = append(d.Orders, askOrder)
				d.OrdersMapping[trade.AskOrder] = askOrder
			}
		}
	}
}
