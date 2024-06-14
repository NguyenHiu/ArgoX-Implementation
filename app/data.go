package app

import (
	"io"
	"math/big"

	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

type VerifyAppData struct {
	Orders []*Order
}

// Encode encodes app data ([]byte) onto an io.Writer.
func (d *VerifyAppData) Encode(w io.Writer) error {
	var encoded_data []byte

	// encode each order
	for _, val := range d.Orders {
		encoded_data = append(encoded_data, val.EncodeOrder()...)
	}

	// write encoded data into writer
	_, err := w.Write(encoded_data)

	return err
}

// A required function of Channel.Data interface
// Clone returns a deep copy of the app data.
func (d *VerifyAppData) Clone() channel.Data {
	_d := *d
	return &_d
}

func (d *VerifyAppData) SendNewOrder(order *Order) {
	d.Orders = append(d.Orders, order)
}

func (d *VerifyAppData) UpdateExistedOrder(orderID uuid.UUID, updatedData OrderUpdatedInfo) {
	for i := 0; i < len(d.Orders); i++ {
		if d.Orders[i].OrderID == orderID {
			if updatedData.Status != "" {
				d.Orders[i].Status = updatedData.Status
			}
			if updatedData.MatchedAmount.Cmp(&big.Int{}) != 0 {
				d.Orders[i].MatchedAmount = updatedData.MatchedAmount
			}
		}
	}
}
