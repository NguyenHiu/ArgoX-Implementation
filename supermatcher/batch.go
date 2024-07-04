package supermatcher

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type MatcherOrder struct {
	Data  *ShadowOrder
	Owner uuid.UUID
}

type Batch struct {
	BatchID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Orders    []*ExpandOrder
	Owner     common.Address
	Signature []byte
}

func (b *Batch) Decode_TransferBatching(_data []byte) error {
	data := bytes.NewBuffer(_data)

	// Batch ID
	if err := binary.Read(data, binary.BigEndian, &b.BatchID); err != nil {
		return err
	}

	// Price
	_price := make([]byte, 32)
	if err := binary.Read(data, binary.BigEndian, &_price); err != nil {
		return err
	}
	b.Price = new(big.Int).SetBytes(_price)

	// Amount
	_amount := make([]byte, 32)
	if err := binary.Read(data, binary.BigEndian, &_amount); err != nil {
		return err
	}
	b.Amount = new(big.Int).SetBytes(_amount)

	// Side
	if err := binary.Read(data, binary.BigEndian, &b.Side); err != nil {
		return err
	}

	// Number of orders
	var _noOrders uint8
	if err := binary.Read(data, binary.BigEndian, &_noOrders); err != nil {
		return err
	}

	for i := 0; i < int(_noOrders); i++ {
		// Shadow order
		shadowOrder := &ShadowOrder{}
		if err := shadowOrder.Decode_TransferBatching(data); err != nil {
			return err
		}

		// Number of executed trades
		var _noExecutedTrades uint8
		if err := binary.Read(data, binary.BigEndian, &_noExecutedTrades); err != nil {
			return err
		}

		// Trades
		executedTrades := []*app.Trade{}
		for j := 0; j < int(_noExecutedTrades); j++ {
			executedTrade := &app.Trade{}
			if err := executedTrade.Decode_TransferBatching(data); err != nil {
				return err
			}
			executedTrades = append(executedTrades, executedTrade)
		}

		// Original Order
		originalOrder := &app.Order{}
		if err := originalOrder.Decode_TransferBatching(data); err != nil {
			return err
		}

		// Append
		b.Orders = append(b.Orders, &ExpandOrder{
			ShadowOrder:   shadowOrder,
			Trades:        executedTrades,
			OriginalOrder: originalOrder,
		})
	}

	// Owner
	_owner := make([]byte, 20)
	if err := binary.Read(data, binary.BigEndian, &_owner); err != nil {
		return err
	}
	b.Owner = common.Address(_owner)

	// Signature
	_signature := make([]byte, 65)
	if err := binary.Read(data, binary.BigEndian, &_signature); err != nil {
		return err
	}
	b.Signature = _signature

	return nil
}

func (b *Batch) Encode_Sign() ([]byte, error) {
	// Get encoded data of all orders
	ordersData := new(bytes.Buffer)
	for _, order := range b.Orders {
		data, err := order.Encode_Sign()
		if err != nil {
			return nil, fmt.Errorf("invalid order data in Sign() func, err: %v", err)
		}
		binary.Write(ordersData, binary.BigEndian, data)
	}

	// Encode packed the batch
	data := new(bytes.Buffer)

	// Batch ID
	if err := binary.Write(data, binary.BigEndian, b.BatchID); err != nil {
		return nil, err
	}

	// Price
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Price)); err != nil {
		return nil, err
	}

	// Amount
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Amount)); err != nil {
		return nil, err
	}

	// Side
	if err := binary.Write(data, binary.BigEndian, b.Side); err != nil {
		return nil, err
	}

	// Number of orders
	if err := binary.Write(data, binary.BigEndian, uint8(len(b.Orders))); err != nil {
		return nil, err
	}

	// Expanded Orders
	if err := binary.Write(data, binary.BigEndian, ordersData.Bytes()); err != nil {
		return nil, err
	}

	// Owner
	if err := binary.Write(data, binary.BigEndian, b.Owner); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (b *Batch) Sign(_prvkey *ecdsa.PrivateKey) error {
	data, err := b.Encode_Sign()
	if err != nil {
		return err
	}

	// Hash the encode packed data
	hasheddata := crypto.Keccak256Hash(data)

	// Sign the batch
	sig, err := crypto.Sign(hasheddata.Bytes(), _prvkey)
	if err != nil {
		log.Fatal(err)
	}
	b.Signature = sig

	return nil
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

func (b *Batch) Equal(_b *Batch) bool {
	if len(b.Orders) != len(_b.Orders) {
		return false
	}
	for idx, order := range b.Orders {
		if !order.Equal(_b.Orders[idx]) {
			return false
		}
	}

	return b.BatchID == _b.BatchID &&
		b.Price.Cmp(_b.Price) == 0 &&
		b.Amount.Cmp(_b.Amount) == 0 &&
		b.Side == _b.Side &&
		b.Owner.Cmp(_b.Owner) == 0 &&
		bytes.Compare(b.Signature, _b.Signature) == 0
}
