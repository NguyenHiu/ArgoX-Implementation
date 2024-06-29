package app

import (
	"encoding/binary"
	"io"
	"sort"

	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

type VerifyAppData struct {
	Orders map[uuid.UUID]*Order
}

/**
 * Encode encodes app data ([]byte) onto an io.Writer.
 * Format: <no_order>(uint64) [<order> <no_msg>(uint64) [<msg>]]
 */
func (d *VerifyAppData) Encode(w io.Writer) error {
	// No orders
	if err := binary.Write(w, binary.BigEndian, uint64(len(d.Orders))); err != nil {
		return err
	}

	ordersKeys := make([]uuid.UUID, 0, len(d.Orders))
	for key := range d.Orders {
		ordersKeys = append(ordersKeys, key)
	}
	sort.Slice(ordersKeys, func(i, j int) bool {
		return ordersKeys[i].String() < ordersKeys[j].String()
	})
	for _, key := range ordersKeys {
		order := d.Orders[key]

		// Order
		if err := binary.Write(w, binary.BigEndian, order.Encode_TransferLightning()); err != nil {
			return err
		}
	}

	return nil
}

// A required function of Channel.Data interface
// Clone returns a deep copy of the app data.
func (d *VerifyAppData) Clone() channel.Data {
	_d := *d

	// Deep copy of the Orders map
	_d.Orders = make(map[uuid.UUID]*Order)
	for key, value := range d.Orders {
		// Assuming Order is a pointer or a simple struct that doesn't require deep copying
		// If Order contains reference types, you would need to further clone those as well
		_d.Orders[key] = value.Clone()
	}
	return &_d
}

func (d *VerifyAppData) SendNewOrder(order *Order) {
	d.Orders[order.OrderID] = order
}

func (d *VerifyAppData) SendNewMessage(message *Message) {
	d.Msgs[message.OrderID] = append(d.Msgs[message.OrderID], message)
}

// TODO:
// TODO:
// TODO: Add component-each-component from mapping-data branch into this branch
// TODO: Note to run the program to check if the new code is right or not!
// TODO:
// TODO:
