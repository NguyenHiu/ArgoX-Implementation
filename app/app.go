package app

import (
	"fmt"
	"io"
	"log"
	"unsafe"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

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
		Orders: []Order{},
	}
}

// DecodeData decodes the channel data.
func (a *VerifyApp) DecodeData(r io.Reader) (channel.Data, error) {
	d := a.InitData()

	// read full data
	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	// order length
	ol := int(unsafe.Sizeof(Order{}))
	// no orders
	no := len(data) / int(ol)
	if int(no) != no {
		log.Fatal("Decode(): decoding suspicious data\n")
	}

	// decode each order
	for i := 0; i < no; i++ {
		order_data := data[i*ol : (i+1)*ol]
		order, err := DecodeOrder(order_data)
		if err != nil {
			log.Fatalf("Decode(): decoding an invalid order, index: %v\n", i)
		}
		d.Orders = append(d.Orders, *order)
	}

	return d, nil
}

// ValidInit checks that the initial state is valid.
func (a *VerifyApp) ValidInit(p *channel.Params, s *channel.State) error {
	if len(p.Parts) != numParts {
		return fmt.Errorf("invalid number of participants: expected %d, got %d", numParts, len(p.Parts))
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
		return fmt.Errorf("Invalid assets: %v", err)
	}

	fromData, ok := from.Data.(*VerifyAppData)
	if !ok {
		panic(fmt.Sprintf("from state: invalid data type: %T", from.Data))
	}

	toData, ok := to.Data.(*VerifyAppData)
	if !ok {
		panic(fmt.Sprintf("to state: invalid data type: %T", from.Data))
	}

	// TODO: checkvlaid transition
	// There are 2 types of transitions:
	// 	1. Trader makes some orders
	//		1.1. If to_state.orders.length > from_state.orders.length
	//		1.2. Check orders from 0 .. from_state.orders.length are not changed
	//		1.3. Check new orders are valid
	//	2. Matcher updates the state of some orders
	//		2.1. If to_state.orders.length == from_state.orders.length
	//		2.2. Check if the changes in orders are valid

	return nil
}
