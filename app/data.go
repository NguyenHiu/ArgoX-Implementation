package app

import (
	"io"

	"perun.network/go-perun/channel"
)

type VerifyAppData struct {
	Orders []Order
}

func (d *VerifyAppData) SendOrder(order Order) {

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
