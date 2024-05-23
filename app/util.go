package app

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"perun.network/go-perun/backend/ethereum/wallet"
)

const (
	ASK           = true
	BID           = false
	STATUS_LENGTH = 1
	numParts      = 2
)

type Order struct {
	Price  float64
	Amount float64
	Side   bool
	Owner  *wallet.Address
	Status string
}

func NewOrder(price, amount float64, side bool, owner *wallet.Address) Order {
	return Order{
		Price:  price,
		Amount: amount,
		Side:   side,
		Owner:  owner,
		Status: "C", // Rename later
	}
}

// Encode Order
// Price > Amount > Side > Owner > Status
func (o *Order) EncodeOrder() []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, o.Price)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, o.Amount)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, o.Side)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, o.Owner.Bytes())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, []byte(o.Status))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

// Decode Order
// Follow the parameter orders when encoding
func DecodeOrder(data []byte) (*Order, error) {
	order := Order{}
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &order.Price)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.LittleEndian, &order.Amount)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.LittleEndian, &order.Side)
	if err != nil {
		return nil, err
	}

	order.Owner = &wallet.Address{}
	err = binary.Read(buf, binary.LittleEndian, order.Owner)
	if err != nil {
		return nil, err
	}

	status_temp := make([]byte, 1)
	err = binary.Read(buf, binary.LittleEndian, &status_temp)
	if err != nil {
		return nil, err
	}
	order.Status = string(status_temp)

	return &order, nil
}
