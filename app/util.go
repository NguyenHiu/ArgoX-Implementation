package app

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
)

const (
	ASK           = true
	BID           = false
	STATUS_LENGTH = 1
	numParts      = 2
)

type Order struct {
	OrderID        uuid.UUID
	Price          float64
	Amount         float64
	Side           bool
	Owner          *wallet.Address
	OwnerSignture  []byte
	Status         string
	MatchedAmoount float64
}

type OrderUpdatedInfo struct {
	Status         string
	MatchedAmoount float64
}

func NewOrder(price, amount float64, side bool, owner *wallet.Address) Order {
	orderId, _ := uuid.NewRandom()
	return Order{
		OrderID:        orderId,
		Price:          price,
		Amount:         amount,
		Side:           side,
		Owner:          owner,
		OwnerSignture:  []byte{},
		Status:         "C", // Replace later
		MatchedAmoount: 0,
	}
}

func (o *Order) Sign(prvkey ecdsa.PrivateKey) error {
	pub, _ := prvkey.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)
	if o.Owner.Cmp(wallet.AsWalletAddr(addr)) != 0 {
		return fmt.Errorf("private key does not match with the order's owner")
	}
	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("invalid uuid")
	}
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, orderID)
	binary.Write(data, binary.LittleEndian, o.Price)
	binary.Write(data, binary.LittleEndian, o.Amount)
	binary.Write(data, binary.LittleEndian, o.Side)
	binary.Write(data, binary.LittleEndian, o.Owner.Bytes())

	hashedData := crypto.Keccak256Hash(data.Bytes())

	sig, err := crypto.Sign(hashedData.Bytes(), &prvkey)
	if err != nil {
		return fmt.Errorf("can not sign the order, err: %v", err)
	}
	o.OwnerSignture = sig

	return nil
}

func (o *Order) IsValidSignature() bool {
	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		return false
	}
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, orderID)
	binary.Write(data, binary.LittleEndian, o.Price)
	binary.Write(data, binary.LittleEndian, o.Amount)
	binary.Write(data, binary.LittleEndian, o.Side)
	binary.Write(data, binary.LittleEndian, o.Owner.Bytes())

	hashedData := crypto.Keccak256Hash(data.Bytes())
	return crypto.VerifySignature(o.Owner.Bytes(), hashedData.Bytes(), o.OwnerSignture)

}

func (o *Order) Equal(_o *Order) bool {
	return (o.OrderID == _o.OrderID &&
		o.Price == _o.Price &&
		o.Amount == _o.Amount &&
		o.Side == _o.Side &&
		o.Owner.Cmp(_o.Owner) == 0 &&
		o.Status == _o.Status &&
		o.MatchedAmoount == _o.MatchedAmoount)
}

// Encode Order
// Price > Amount > Side > Owner > Status
func (o *Order) EncodeOrder() []byte {
	buf := new(bytes.Buffer)

	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		fmt.Println("invalid uuid")
	}
	err = binary.Write(buf, binary.LittleEndian, orderID)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, o.Price)
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

	err = binary.Write(buf, binary.LittleEndian, o.OwnerSignture)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, []byte(o.Status))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.LittleEndian, o.MatchedAmoount)
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

	orderIDTemp := make([]byte, 16)
	err := binary.Read(buf, binary.LittleEndian, &orderIDTemp)
	if err != nil {
		return nil, err
	}
	err = order.OrderID.UnmarshalBinary(orderIDTemp)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.LittleEndian, &order.Price)
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

	ownerSign := make([]byte, 65)
	err = binary.Read(buf, binary.LittleEndian, &ownerSign)
	if err != nil {
		return nil, err
	}
	order.OwnerSignture = ownerSign

	status_temp := make([]byte, 1)
	err = binary.Read(buf, binary.LittleEndian, &status_temp)
	if err != nil {
		return nil, err
	}
	order.Status = string(status_temp)

	err = binary.Read(buf, binary.LittleEndian, &order.MatchedAmoount)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
