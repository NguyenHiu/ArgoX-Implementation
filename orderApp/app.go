package orderApp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

var _logger = logger.NewLogger("orderApp", logger.None, logger.None)

type OrderApp struct {
	Addr wallet.Address
}

func NewOrderApp(addr wallet.Address) *OrderApp {
	return &OrderApp{
		Addr: addr,
	}
}

func (a *OrderApp) Def() wallet.Address {
	return a.Addr
}

func (a *OrderApp) InitData() *OrderAppData {
	return &OrderAppData{
		Orders:        make([]*Order, 0),
		OrdersMapping: make(map[uuid.UUID]*Order),
	}
}

/**
 * DecodeData decodes the channel data.
 */
func (a *OrderApp) DecodeData(r io.Reader) (channel.Data, error) {
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
	return d, nil
}

// ValidInit checks that the initial state is valid.
func (a *OrderApp) ValidInit(p *channel.Params, s *channel.State) error {
	if len(p.Parts) != constants.NUM_PARTS {
		return fmt.Errorf("invalid number of participants: expected %d, got %d", constants.NUM_PARTS, len(p.Parts))
	}

	appData, ok := s.Data.(*OrderAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", s.Data)
	}

	if len(appData.Orders) != 0 {
		return fmt.Errorf("invalid starting")
	}

	return nil
}

// ValidTransition is called whenever the channel state transitions.
func (a *OrderApp) ValidTransition(params *channel.Params, from, to *channel.State, idx channel.Index) error {
	err := channel.AssetsAssertEqual(from.Assets, to.Assets)
	if err != nil {
		_logger.Error("invalid assets: %v\n", err)
		return fmt.Errorf("invalid assets: %v", err)
	}

	// Get data
	fromData, ok := from.Data.(*OrderAppData)
	if !ok {
		_logger.Error("from state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("from state: invalid data type: %T", from.Data)
	}

	toData, ok := to.Data.(*OrderAppData)
	if !ok {
		_logger.Error("to state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("to state: invalid data type: %T", from.Data)
	}

	// Check change
	if len(toData.Orders) < len(fromData.Orders) {
		_logger.Error("invalid transition: the number of orders in new state is incorrect\n")
		return fmt.Errorf("invalid transition: the number of orders in new state is incorrect")
	}

	// Check change detail
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
			if !v.Equal(_v) {
				_logger.Error("invalid transition: \n")
				return fmt.Errorf("invalid transition: ")
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

func (a *OrderApp) SendNewOrders(s *channel.State, orders []*Order) error {
	d, ok := s.Data.(*OrderAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	d.SendNewOrders(orders)

	for _, order := range orders {
		if order.OrderID == EndID {
			_logger.Debug("Got final order\n")
			s.IsFinal = true
			s.Balances = d.computeFinalBalances(s.Balances)
		}
	}

	return nil
}
