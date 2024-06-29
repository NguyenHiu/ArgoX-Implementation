package app

import (
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
		Orders: make(map[uuid.UUID]*Order),
	}
}

/**
 * DecodeData decodes the channel data.
 * Format: <no_order>(uint64) [<order> <no_msg>(uint64) [<msg>]]
 */
func (a *VerifyApp) DecodeData(r io.Reader) (channel.Data, error) {
	d := a.InitData()

	// Read data
	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	// Get no orders
	noOrders := int(binary.BigEndian.Uint64(data[:8]))
	from := 8
	for i := 0; i < noOrders; i++ {
		// Get order
		order, err := Order_Decode_TransferLightning(data[from : from+constants.LIGHTNING_ORDER_SIZE])
		if err != nil {
			_logger.Error("Order Decode Transfer Lightning fail, err: %v\n", err)
			return nil, err
		}
		from += constants.LIGHTNING_ORDER_SIZE

		d.Orders[order.OrderID] = order
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

	return nil
}

// ValidTransition is called whenever the channel state transitions.
func (a *VerifyApp) ValidTransition(params *channel.Params, from, to *channel.State, idx channel.Index) error {
	err := channel.AssetsAssertEqual(from.Assets, to.Assets)
	if err != nil {
		_logger.Error("invalid assets: %v\n", err)
		return fmt.Errorf("invalid assets: %v", err)
	}

	// Get data
	fromData, ok := from.Data.(*VerifyAppData)
	if !ok {
		_logger.Error("from state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("from state: invalid data type: %T", from.Data)
	}

	toData, ok := to.Data.(*VerifyAppData)
	if !ok {
		_logger.Error("to state: invalid data type: %T\n", from.Data)
		return fmt.Errorf("to state: invalid data type: %T", from.Data)
	}

	// Check change
	if len(fromData.Orders)+1 != len(toData.Orders) && len(fromData.Orders) != len(toData.Orders) {
		_logger.Error("invalid transition: the number of orders in new state is incorrect\n")
		return fmt.Errorf("invalid transition: the number of orders in new state is incorrect")
	}

	// Check change detail
	flag := false
	for k, v := range toData.Orders {
		_v, ok := fromData.Orders[k]
		if !ok {
			if flag {
				_logger.Error("invalid transition: too much orders for a state transition\n")
				return fmt.Errorf("invalid transition: too much orders for a state transition")
			}
			// Validate new order
			if !v.IsValidSignature() {
				_logger.Error("invalid transition: the new order is not valid\n")
				return fmt.Errorf("invalid transition: the new order is not valid")
			}
			flag = true
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

func (a *VerifyApp) SendNewOrder(s *channel.State, order *Order) error {
	d, ok := s.Data.(*VerifyAppData)
	if !ok {
		return fmt.Errorf("invalid data type: %T", d)
	}

	d.SendNewOrder(order)

	if order.OrderID == EndID {
		s.IsFinal = true
		s.Balances = d.computeFinalBalances(s.Balances)
	}
	return nil
}
