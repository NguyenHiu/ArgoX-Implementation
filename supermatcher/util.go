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

func (o *Order) Encode_TranferBatching() ([]byte, error) {
	data := new(bytes.Buffer)
	id, err := o.OrderID.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, id); err != nil {
		fmt.Println("Can not convert order id to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(o.Price)); err != nil {
		fmt.Println("Can not convert order price to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(o.Amount)); err != nil {
		fmt.Println("Can not convert order amount to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, o.Side); err != nil {
		fmt.Println("Can not convert order side to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, o.Owner.Bytes()); err != nil {
		fmt.Println("Can not convert order owner to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, o.Signature); err != nil {
		fmt.Println("Can not convert order signature to []byte")
		return []byte{}, err
	}

	return data.Bytes(), nil
}

func (o *Order) Encode_Sign() ([]byte, error) {
	data := new(bytes.Buffer)
	id, err := o.OrderID.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, id); err != nil {
		fmt.Println("Can not convert order id to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(o.Price)); err != nil {
		fmt.Println("Can not convert order price to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(o.Amount)); err != nil {
		fmt.Println("Can not convert order amount to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, o.Side); err != nil {
		fmt.Println("Can not convert order side to []byte")
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, o.Owner.Bytes()); err != nil {
		fmt.Println("Can not convert order owner to []byte")
		return []byte{}, err
	}

	hashedData := crypto.Keccak256Hash(data.Bytes())

	return hashedData.Bytes(), nil
}

func Order_Decode_TransferBatching(data []byte) (*Order, error) {
	order := Order{}
	buf := bytes.NewBuffer(data)

	id := make([]byte, 16)
	if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
		return nil, err
	}
	if err := order.OrderID.UnmarshalBinary(id); err != nil {
		return nil, err
	}

	price := make([]byte, 32)
	if err := binary.Read(buf, binary.BigEndian, &price); err != nil {
		return nil, err
	}
	order.Price = new(big.Int).SetBytes(price)

	amount := make([]byte, 32)
	if err := binary.Read(buf, binary.BigEndian, &amount); err != nil {
		return nil, err
	}
	order.Amount = new(big.Int).SetBytes(amount)

	if err := binary.Read(buf, binary.BigEndian, &order.Side); err != nil {
		return nil, err
	}

	if err := binary.Read(buf, binary.BigEndian, &order.Owner); err != nil {
		return nil, err
	}

	sign := make([]byte, 65)
	if err := binary.Read(buf, binary.BigEndian, &sign); err != nil {
		return nil, err
	}
	order.Signature = sign

	return &order, nil
}

func (o *Order) IsValidSignature() bool {
	encodedOrder, err := o.Encode_Sign()
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

func Batch_Decode_TransferBatching(data []byte) (*Batch, error) {
	batch := Batch{}
	buf := bytes.NewBuffer(data)

	id := make([]byte, 16)
	if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
		return nil, err
	}
	if err := batch.BatchID.UnmarshalBinary(id); err != nil {
		return nil, err
	}

	price := make([]byte, 32)
	if err := binary.Read(buf, binary.BigEndian, &price); err != nil {
		return nil, err
	}
	batch.Price = new(big.Int).SetBytes(price)

	amount := make([]byte, 32)
	if err := binary.Read(buf, binary.BigEndian, &amount); err != nil {
		return nil, err
	}
	batch.Amount = new(big.Int).SetBytes(amount)

	if err := binary.Read(buf, binary.BigEndian, &batch.Side); err != nil {
		return nil, err
	}

	var orderLength uint8
	if err := binary.Read(buf, binary.BigEndian, &orderLength); err != nil {
		return nil, err
	}

	orders := []*Order{}
	orderData := make([]byte, 166*orderLength)
	if err := binary.Read(buf, binary.BigEndian, orderData); err != nil {
		return nil, err
	}
	for i := 0; i < int(orderLength); i++ {
		order, err := Order_Decode_TransferBatching(orderData[(i * 166):((i + 1) * 166)])
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	batch.Orders = orders

	if err := binary.Read(buf, binary.BigEndian, &batch.Owner); err != nil {
		return nil, err
	}

	sign := make([]byte, 65)
	if err := binary.Read(buf, binary.BigEndian, sign); err != nil {
		return nil, err
	}

	batch.Signature = sign
	return &batch, nil
}

func (b *Batch) Encode_Sign() ([]byte, error) {

	ordersData := new(bytes.Buffer)
	for _, order := range b.Orders {
		data, err := order.Encode_TranferBatching()
		if err != nil {
			return []byte{}, fmt.Errorf("invalid order data in Sign() func, err: %v", err)
		}
		binary.Write(ordersData, binary.BigEndian, data)
	}

	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, b.BatchID)
	binary.Write(data, binary.BigEndian, PaddingToUint256(b.Price))
	binary.Write(data, binary.BigEndian, PaddingToUint256(b.Amount))
	binary.Write(data, binary.BigEndian, b.Side)
	binary.Write(data, binary.BigEndian, uint8(len(b.Orders)))
	binary.Write(data, binary.BigEndian, ordersData.Bytes())
	binary.Write(data, binary.BigEndian, b.Owner)

	return data.Bytes(), nil
}

func (b *Batch) IsValidSignature() bool {
	encodedBatch, err := b.Encode_Sign()
	if err != nil {
		return false
	}

	hasheddata := crypto.Keccak256Hash(encodedBatch)

	pubkey, err := crypto.SigToPub(hasheddata.Bytes(), b.Signature)
	if err != nil {
		return false
	}

	if crypto.PubkeyToAddress(*pubkey).Cmp(b.Owner) != 0 {
		return false
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(pubkey), hasheddata.Bytes(), b.Signature[:64])
}

func PaddingToUint256(num *big.Int) []byte {
	return append(make([]byte, 32-len(num.Bytes())), num.Bytes()...)
}
