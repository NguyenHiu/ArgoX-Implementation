package app

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
)

const (
	ASK           = true
	BID           = false
	STATUS_LENGTH = 1
	numParts      = 2
	ORDER_SIZE    = 127
)

type Order struct {
	OrderID        uuid.UUID
	Price          int64
	Amount         int64
	Side           bool
	Owner          *wallet.Address
	OwnerSignture  []byte
	Status         string
	MatchedAmoount int64
}

type OrderUpdatedInfo struct {
	Status         string
	MatchedAmoount int64
}

// The `status` parameter should be "P" at the init phase,
// but allowing the `status` parameter to be passed is for testing purposes.
func NewOrder(price, amount int64, side bool, owner *wallet.Address, status string) Order {
	orderId, _ := uuid.NewRandom()
	return Order{
		OrderID:        orderId,
		Price:          price,
		Amount:         amount,
		Side:           side,
		Owner:          owner,
		OwnerSignture:  []byte{},
		Status:         status, // Replace later
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

	pub, err := crypto.SigToPub(hashedData.Bytes(), o.OwnerSignture)
	if err != nil {
		fmt.Printf("Cannot recover public key from signature, error: %v\n", err)
		return false
	}
	_owner := wallet.AsWalletAddr(crypto.PubkeyToAddress(*pub))
	if _owner.Cmp(o.Owner) != 0 {
		fmt.Println("Provided public key does not match with the order's owner")
		return false
	}
	pubBytes := crypto.FromECDSAPub(pub)
	return crypto.VerifySignature(pubBytes, hashedData.Bytes(), o.OwnerSignture[:len(o.OwnerSignture)-1])

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

func (d *VerifyAppData) CheckFinal() bool {
	l := len(d.Orders)
	return l != 0 && d.Orders[l-1].Status == "F"
}

func computeFinalBalances(orders []*Order, bals channel.Balances) channel.Balances {
	matcherReceivedAmount := int64(0)

	for i := 0; i < len(orders); i++ {
		// if orders[i].Status == "M" {
		if orders[i].Status != "F" {
			if !orders[i].Side {
				matcherReceivedAmount += orders[i].Price
			} else {
				matcherReceivedAmount -= orders[i].Price
			}
		}
		// }
	}

	fmt.Printf("matcherReceivedAmount: %v\n", matcherReceivedAmount)

	finalBals := bals.Clone()
	for i := range finalBals {
		bigIntAmount := big.NewInt(matcherReceivedAmount)
		finalBals[i][0] = new(big.Int).Sub(bals[i][0], bigIntAmount)
		finalBals[i][1] = new(big.Int).Add(bals[i][1], bigIntAmount)
	}

	return finalBals
}
