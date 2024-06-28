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
	Msgs   map[uuid.UUID][]*Message
	Test   []int
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

		// No msgs
		if err := binary.Write(w, binary.BigEndian, uint64(len(d.Msgs[order.OrderID]))); err != nil {
			return err
		}
		// Message
		for _, m := range d.Msgs[order.OrderID] {
			if err := binary.Write(w, binary.BigEndian, m.Encode_TransferLightning()); err != nil {
				return err
			}
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

	// Deep copy of the Msgs map
	_d.Msgs = make(map[uuid.UUID][]*Message)
	for key, slice := range d.Msgs {
		copiedSlice := make([]*Message, len(slice))
		for _idx, _msg := range slice {
			copiedSlice[_idx] = _msg.Clone()
		}
		_d.Msgs[key] = copiedSlice
	}
	return &_d
}

func (d *VerifyAppData) SendNewOrder(order *Order) {
	// _logger.Debug("SendNewOrder, [1] len(d.Orders): %v\n", len(d.Orders))
	d.Orders[order.OrderID] = order
	// _logger.Debug("SendNewOrder, [2] len(d.Orders): %v\n", len(d.Orders))
	d.Msgs[order.OrderID] = []*Message{}
}

// TODO: Instead of updating existing orders, create new message to notify the match event
func (d *VerifyAppData) SendMessage(message *Message) {
	d.Msgs[message.OrderID] = append(d.Msgs[message.OrderID], message)
}
