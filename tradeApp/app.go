package tradeApp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

var _logger = logger.NewLogger("tradeApp", logger.None, logger.None)

type TradeApp struct {
	Addr wallet.Address
}

func NewTradeApp(addr wallet.Address) *TradeApp {
	return &TradeApp{
		Addr: addr,
	}
}

func (a *TradeApp) Def() wallet.Address {
	return a.Addr
}

func (a *TradeApp) InitData() *TradeAppData {
	return &TradeAppData{
		Orders:        make([]*Order, 0),
		OrdersMapping: make(map[uuid.UUID]*Order),
		Trades:        make([]*Trade, 0),
		TradesMapping: make(map[uuid.UUID]*Trade),
		BidToTrade:    make(map[uuid.UUID][]*Trade),
		AskToTrade:    make(map[uuid.UUID][]*Trade),
	}
}

/**
 * DecodeData decodes the channel data.
 * Format: <no_order> <order>... <no message list> [<no message> <message>...]...
 */
func (a *TradeApp) DecodeData(r io.Reader) (channel.Data, error) {
	d := a.InitData()

	// Read data
	_data, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	data := bytes.NewBuffer(_data)

	// Get no orders
	var noOrders uint8
	if err := binary.Read(data, binary.BigEndian, &noOrders); err != nil {
		return nil, err
	}

	for i := 0; i < int(noOrders); i++ {
		// Get order
		order := &Order{}
		if err := order.Decode_TransferLightning(data); err != nil {
			_logger.Error("Order Decode Transfer Lightning fail, err: %v\n", err)
			return nil, err
		}
		// Store order
		d.Orders = append(d.Orders, order)
		d.OrdersMapping[order.OrderID] = order
	}

	// No Trades
	var noTrades uint8
	if err := binary.Read(data, binary.BigEndian, &noTrades); err != nil {
		return nil, err
	}

	// Each Trade
	for i := 0; i < int(noTrades); i++ {
		trade := &Trade{}
		if err := trade.Decode_TransferLightning(data); err != nil {
			return nil, err
		}
		d.BidToTrade[trade.BidOrder] = append(d.BidToTrade[trade.BidOrder], trade)
		d.AskToTrade[trade.AskOrder] = append(d.AskToTrade[trade.AskOrder], trade)
		d.Trades = append(d.Trades, trade)
		d.TradesMapping[trade.TradeID] = trade
	}

	return d, nil
}

// ValidInit checks that the initial state is valid.
func (a *TradeApp) ValidInit(p *channel.Params, s *channel.State) error {
	if len(p.Parts) != constants.NUM_PARTS {
		return fmt.Errorf("invalid number of participants: expected %d, got %d", constants.NUM_PARTS, len(p.Parts))
	}

	appData, ok := s.Data.(*TradeAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", s.Data)
	}

	if len(appData.Orders) != 0 {
		return fmt.Errorf("invalid starting")
	}

	if len(appData.Trades) != 0 {
		return fmt.Errorf("invalid starting")
	}

	if len(appData.AskToTrade) != 0 {
		return fmt.Errorf("invalid starting")
	}

	if len(appData.BidToTrade) != 0 {
		return fmt.Errorf("invalid starting")
	}

	return nil
}

// ValidTransition is called whenever the channel state transitions.
func (a *TradeApp) ValidTransition(params *channel.Params, from, to *channel.State, idx channel.Index) error {
	err := channel.AssetsAssertEqual(from.Assets, to.Assets)
	if err != nil {
		_logger.Error("invalid assets: %v\n", err)
		return fmt.Errorf("invalid assets: %v", err)
	}

	// Get data
	fromData, ok := from.Data.(*TradeAppData)
	if !ok {
		_logger.Error("from state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("from state: invalid data type: %T", from.Data)
	}

	toData, ok := to.Data.(*TradeAppData)
	if !ok {
		_logger.Error("to state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("to state: invalid data type: %T", from.Data)
	}

	// Check change
	if len(toData.Orders) < len(fromData.Orders) {
		_logger.Error("invalid transition: the number of orders in new state is incorrect\n")
		return fmt.Errorf("invalid transition: the number of orders in new state is incorrect")
	}

	// Check Orders
	for _, v := range toData.Orders {
		_v, ok := fromData.OrdersMapping[v.OrderID]
		if !ok {
			// Validate new order
			if !v.IsValidSignature() {
				_logger.Error("invalid transition: the new order is not valid\n")
				return fmt.Errorf("invalid transition: the new order is not valid")
			}
		} else {
			// Check if the order stays the same
			// _logger.Debug("v: %v\n", v)
			// _logger.Debug("_v: %v\n", _v)
			if !v.Equal(_v) {
				_logger.Error("invalid transition: \n")
				return fmt.Errorf("invalid transition: ")
			}
		}
	}

	// Check Trades
	for _, v := range toData.Trades {
		_v, ok := fromData.TradesMapping[v.TradeID]
		if !ok {
			// Validate new trade
			if !v.IsValidSignature() {
				_logger.Error("invalid transition: the new trade is not valid\n")
				return fmt.Errorf("invalid transition: the new trade is not valid")
			}
		} else {
			// Check if the old trades stay the same
			if !v.Equal(_v) {
				_logger.Debug("v: %v\n", v)
				_logger.Debug("_v: %v\n", _v)
				_logger.Error("invalid transition: the old trades were changed\n")
				return fmt.Errorf("invalid transition: the old trades were changed")
			}

			// Check if the matched amount is valid
			_, ok1 := toData.OrdersMapping[v.BidOrder]
			_, ok2 := toData.OrdersMapping[v.AskOrder]
			if !ok1 && !ok2 {
				_logger.Error("invalid transition: trade invalid\n")
				return fmt.Errorf("invalid transition: trade invalid")
			}
			if ok1 {
				amount := new(big.Int).Set(toData.OrdersMapping[v.BidOrder].Amount)
				// _logger.Debug("amount: %v\n", amount)
				for _, trade := range toData.BidToTrade[v.BidOrder] {
					amount = new(big.Int).Sub(amount, trade.Amount)
				}
				if amount.Cmp(new(big.Int)) == -1 {
					_logger.Debug("bid order's id: %v\n", v.BidOrder.String())
					_logger.Debug("original amount: %v\n", toData.OrdersMapping[v.BidOrder].Amount)
					_logger.Debug("amount: %v\n", amount)
					_logger.Debug("v.Amount: %v\n", v.Amount)
					_logger.Error("invalid transition: trade's amount invalid\n")
					return fmt.Errorf("invalid transition: trade's amount invalid")
				}
			}
			if ok2 {
				amount := new(big.Int).Set(toData.OrdersMapping[v.AskOrder].Amount)
				for _, trade := range toData.BidToTrade[v.AskOrder] {
					amount = new(big.Int).Sub(amount, trade.Amount)
				}
				if amount.Cmp(new(big.Int)) == -1 {
					_logger.Debug("ask order's id: %v\n", v.AskOrder.String())
					_logger.Debug("original amount: %v\n", toData.OrdersMapping[v.AskOrder].Amount)
					_logger.Debug("amount: %v\n", amount)
					_logger.Debug("v.Amount: %v\n", v.Amount)
					_logger.Error("invalid transition: trade's amount invalid\n")
					return fmt.Errorf("invalid transition: trade's amount invalid")
				}
			}

			// Check signature
			if !v.IsValidSignature() {
				_logger.Error("invalid transition: the new trade's signature is invalid\n")
				return fmt.Errorf("invalid transition: the new trade's signature is invalid")
			}
		}
	}

	isFinal := toData.CheckFinal()
	if isFinal != to.IsFinal {
		return fmt.Errorf("final flag: expected %v, got %v", to.IsFinal, isFinal)
	}

	expectedAllocation := from.Allocation.Clone()
	if isFinal {
		expectedAllocation.Balances = toData.computeFinalBalances(from.Balances)
	}
	if err := expectedAllocation.Equal(&to.Allocation); err != nil {
		return errors.WithMessagef(err, "wrong allocation: expected %v, got %v", expectedAllocation, to.Allocation)
	}

	return nil
}

func (a *TradeApp) SendNewTrades(s *channel.State, trades []*Trade, bidOrder, askOrder *Order, isBid bool) error {
	d, ok := s.Data.(*TradeAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	d.SendNewTrades(trades, bidOrder, askOrder, isBid)

	for _, trade := range trades {
		if trade.TradeID == EndID {
			s.IsFinal = true
			s.Balances = d.computeFinalBalances(s.Balances)
		}
	}

	return nil
}
