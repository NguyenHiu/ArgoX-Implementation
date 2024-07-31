package app

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

var _logger = logger.NewLogger("app", logger.None, logger.None)

type VerifyApp struct {
	Addr wallet.Address
}

func NewVerifyApp(addr wallet.Address) *VerifyApp {
	return &VerifyApp{
		Addr: addr,
	}
}

func (a *VerifyApp) Def() wallet.Address {
	return a.Addr
}

func (a *VerifyApp) InitData() *VerifyAppData {
	return &VerifyAppData{
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
func (a *VerifyApp) DecodeData(r io.Reader) (channel.Data, error) {
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
			//IMHERETODEBUG_logger.Error("Order Decode Transfer Lightning fail, err: %v\n", err)
			return nil, err
		}
		// Store order
		d.Orders = append(d.Orders, order)
		d.OrdersMapping[order.OrderID] = order
	}

	// No Trades
	var noTrades uint8
	if err := binary.Read(data, binary.BigEndian, &noOrders); err != nil {
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
func (a *VerifyApp) ValidInit(p *channel.Params, s *channel.State) error {
	if len(p.Parts) != constants.NUM_PARTS {
		return fmt.Errorf("invalid number of participants: expected %d, got %d", constants.NUM_PARTS, len(p.Parts))
	}

	appData, ok := s.Data.(*VerifyAppData)
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
func (a *VerifyApp) ValidTransition(params *channel.Params, from, to *channel.State, idx channel.Index) error {
	err := channel.AssetsAssertEqual(from.Assets, to.Assets)
	if err != nil {
		//IMHERETODEBUG_logger.Error("invalid assets: %v\n", err)
		return fmt.Errorf("invalid assets: %v", err)
	}

	// Get data
	fromData, ok := from.Data.(*VerifyAppData)
	if !ok {
		//IMHERETODEBUG_logger.Error("from state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("from state: invalid data type: %T", from.Data)
	}

	toData, ok := to.Data.(*VerifyAppData)
	if !ok {
		//IMHERETODEBUG_logger.Error("to state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("to state: invalid data type: %T", from.Data)
	}

	// Check change
	if len(toData.Orders) < len(fromData.Orders) {
		//IMHERETODEBUG_logger.Error("invalid transition: the number of orders in new state is incorrect\n")
		return fmt.Errorf("invalid transition: the number of orders in new state is incorrect")
	}

	// Check change detail
	for _, v := range toData.Orders {
		_v, ok := fromData.OrdersMapping[v.OrderID]
		if !ok {
			// Validate new order
			if !v.IsValidSignature() {
				//IMHERETODEBUG_logger.Error("invalid transition: the new order is not valid\n")
				return fmt.Errorf("invalid transition: the new order is not valid")
			}
		} else {
			// Check if the order stays the same
			if !v.Equal(_v) {
				//IMHERETODEBUG_logger.Error("invalid transition: \n")
				return fmt.Errorf("invalid transition: ")
			}

			total := new(big.Int)
			if v.Side == constants.BID {
				for _, trade := range toData.BidToTrade[v.OrderID] {
					total = new(big.Int).Add(total, trade.Amount)
				}
			} else {
				for _, trade := range toData.AskToTrade[v.OrderID] {
					total = new(big.Int).Add(total, trade.Amount)
				}
			}
			// //IMHERETODEBUG_logger.Debug("total: %v\n", total)
			// //IMHERETODEBUG_logger.Debug("v.Amount: %v\n", v.Amount)
			if total.Cmp(v.Amount) == 1 {
				return fmt.Errorf("invalid transition: the amount accumulate is over the order")
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

func (a *VerifyApp) SendNewOrders(s *channel.State, orders []*Order) error {
	d, ok := s.Data.(*VerifyAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	d.SendNewOrders(orders)

	for _, order := range orders {
		if order.OrderID == EndID {
			//IMHERETODEBUG_logger.Debug("Got final order\n")
			s.IsFinal = true
			s.Balances = d.computeFinalBalances(s.Balances)
		}
	}

	return nil
}

func (a *VerifyApp) SendNewTrades(s *channel.State, trades []*Trade) error {
	d, ok := s.Data.(*VerifyAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	d.SendNewTrades(trades)
	return nil
}
