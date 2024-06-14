package supermatcher

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

/////// ORDER

type Order struct {
	OrderID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Owner     common.Address
	Signature []byte
}

func (o *Order) Equal(_o *Order) bool {
	return o.OrderID == _o.OrderID &&
		o.Price == _o.Price &&
		o.Amount == _o.Amount &&
		o.Side == _o.Side &&
		o.Owner.Cmp(_o.Owner) == 0 &&
		bytes.Equal(o.Signature, _o.Signature)
}

func (o *Order) Encode() ([]byte, error) {
	data := new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, o.OrderID)
	if err != nil {
		fmt.Println("Can not convert order id to []byte")
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, o.Price)
	if err != nil {
		fmt.Println("Can not convert order price to []byte")
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, o.Amount)
	if err != nil {
		fmt.Println("Can not convert order amount to []byte")
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, o.Side)
	if err != nil {
		fmt.Println("Can not convert order side to []byte")
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, o.Owner.Bytes())
	if err != nil {
		fmt.Println("Can not convert order owner to []byte")
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, o.Signature)
	if err != nil {
		fmt.Println("Can not convert order signature to []byte")
		return []byte{}, err
	}

	return data.Bytes(), nil
}

func (o *Order) IsValidSignature() bool {
	encodedOrder, err := o.Encode()
	if err != nil {
		return false
	}

	pubkey, err := crypto.SigToPub(encodedOrder, o.Signature)
	if err != nil {
		return false
	}

	if crypto.PubkeyToAddress(*pubkey).Cmp(o.Owner) != 0 {
		return false
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(pubkey), encodedOrder, o.Signature[:64])
}

/////// BATCH

type Batch struct {
	BatchID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Orders    []*Order
	Owner     common.Address
	Signature []byte
}

func (b *Batch) Equal(_b *Batch) bool {
	if b.BatchID != _b.BatchID ||
		b.Price != _b.Price ||
		b.Amount != _b.Amount ||
		b.Side != _b.Side ||
		b.Owner.Cmp(_b.Owner) != 0 ||
		bytes.Equal(b.Signature, _b.Signature) ||
		len(b.Orders) != len(_b.Orders) {
		return false
	}

	for idx, order := range b.Orders {
		if !order.Equal(_b.Orders[idx]) {
			return false
		}
	}

	return true
}

func (b *Batch) encode() ([]byte, error) {
	data := new(bytes.Buffer)
	if err := binary.Write(data, binary.BigEndian, b.BatchID); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Price); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Amount); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Side); err != nil {
		return []byte{}, err
	}

	ordersbyte := new(bytes.Buffer)
	for _, order := range b.Orders {
		encodedOrder, err := order.Encode()
		if err != nil {
			return []byte{}, err
		}

		if err := binary.Write(ordersbyte, binary.BigEndian, encodedOrder); err != nil {
			return []byte{}, err
		}
	}

	if err := binary.Write(data, binary.BigEndian, ordersbyte); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Owner.Bytes()); err != nil {
		return []byte{}, err
	}

	return data.Bytes(), nil
}

func (b *Batch) IsValidSignature() bool {
	encodedBatch, err := b.encode()
	if err != nil {
		return false
	}

	pubkey, err := crypto.SigToPub(encodedBatch, b.Signature)
	if err != nil {
		return false
	}

	if crypto.PubkeyToAddress(*pubkey).Cmp(b.Owner) != 0 {
		return false
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(pubkey), encodedBatch, b.Signature[:64])
}
