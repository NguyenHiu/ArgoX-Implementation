package orderApp

import (
	"encoding/binary"
	"io"

	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

type OrderAppData struct {
	Orders        []*Order
	OrdersMapping map[uuid.UUID]*Order
}

/**
 * Encode encodes app data ([]byte) onto an io.Writer.
 */
func (d *OrderAppData) Encode(w io.Writer) error {
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

	return nil
}

func (d *OrderAppData) Clone() channel.Data {
	cloned := &OrderAppData{
		Orders:        make([]*Order, len(d.Orders)),
		OrdersMapping: make(map[uuid.UUID]*Order),
	}

	// Clone Orders
	for i, order := range d.Orders {
		cloned.Orders[i] = order.Clone() // Assuming Order has a Clone method
	}

	// Clone OrdersMapping
	for key, order := range d.OrdersMapping {
		cloned.OrdersMapping[key] = order.Clone() // Assuming Order has a Clone method
	}

	return cloned
}

func (d *OrderAppData) SendNewOrders(orders []*Order) {
	for _, order := range orders {
		_, ok := d.OrdersMapping[order.OrderID]
		if !ok {
			d.Orders = append(d.Orders, order)
			d.OrdersMapping[order.OrderID] = order
		}
	}
}
