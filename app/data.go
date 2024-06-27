package app

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

type VerifyAppData struct {
	// Orders     []*Order
	Orders map[uuid.UUID]*Order
	Msgs   map[uuid.UUID][]*Message
}

// GUESS: Maybe dump data from d into the w to be validated
// It returns erorr only if I w.write()
// Encode encodes app data ([]byte) onto an io.Writer.
func (d *VerifyAppData) Encode(w io.Writer) error {
	fmt.Println("actually")
	encodedData := new(bytes.Buffer)

	for _, v := range d.Orders {
		binary.Write(encodedData, binary.BigEndian, v.Encode_TransferLightning())
	}

	for _, v := range d.Msgs {
		for _, _v := range v {
			binary.Write(encodedData, binary.BigEndian, _v.Encode_TransferLightning())
		}
	}

	// // No orders
	// binary.Write(encodedData, binary.BigEndian, uint64(len(d.Orders)))
	// _logger.Debug("len(d.Orders): %v\n", uint64(len(d.Orders)))
	// for _, v := range d.Orders {
	// 	// Order
	// 	binary.Write(encodedData, binary.BigEndian, v.Encode_TransferLightning())
	// 	// _logger.Debug("lv.Encode_TransferLightning(): %v\n", v.Encode_TransferLightning())

	// 	// No msgs
	// 	binary.Write(encodedData, binary.BigEndian, uint64(len(d.Msgs[v.OrderID])))
	// 	// _logger.Debug("len(d.Msgs[v.OrderID]): %v\n", uint64(len(d.Msgs[v.OrderID])))
	// 	// Message
	// 	for _, m := range d.Msgs[v.OrderID] {
	// 		binary.Write(encodedData, binary.BigEndian, m.Encode_TransferLightning())
	// 		// _logger.Debug("m.Encode_TransferLightning(): %v\n", m.Encode_TransferLightning())
	// 	}
	// }

	// Write encoded data into writer
	_, _ = w.Write(encodedData.Bytes())
	// if err != nil {

	// 	return err
	// }

	// _logger.Debug("%v-%v\n", len(d.Orders), len(d.Msgs))
	// _logger.Debug("\n")
	// _logger.Debug("len of Orders: %v\n", len(d.Orders))
	// _logger.Debug("len of Msg: %v\n", len(d.Msgs))
	// _logger.Debug("BRUH IM AT THE BOTTOM OF THE ENCODE FUNCTION\n")
	// _logger.Debug("encoded data: %v\n", encodedData.Bytes())
	return nil
}

// A required function of Channel.Data interface
// Clone returns a deep copy of the app data.
func (d *VerifyAppData) Clone() channel.Data {
	_d := *d
	_d.Orders = make(map[uuid.UUID]*Order)
	for key, value := range d.Orders {
		_d.Orders[key] = value.Clone()
	}
	_d.Msgs = make(map[uuid.UUID][]*Message)
	for key, slice := range d.Msgs {
		copiedSlice := make([]*Message, len(slice))
		copy(copiedSlice, slice)
		_d.Msgs[key] = copiedSlice
	}
	return &_d
}

func (d *VerifyAppData) SendNewOrder(order *Order) {
	// d.Orders = append(d.Orders, order)
	d.Orders[order.OrderID] = order
	d.Msgs[order.OrderID] = []*Message{}
}

func (d *VerifyAppData) UpdateExistedOrder(orderID uuid.UUID, updatedData OrderUpdatedInfo) {
	// for i := 0; i < len(d.Orders); i++ {
	// 	if d.Orders[i].OrderID == orderID {
	// 		if updatedData.Status != "" {
	// 			d.Orders[i].Status = updatedData.Status
	// 		}
	// 		if updatedData.MatchedAmount.Cmp(&big.Int{}) != 0 {
	// 			d.Orders[i].MatchedAmount = updatedData.MatchedAmount
	// 		}
	// 	}
	// }
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
