package supermatcher

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/tradeApp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ExpandOrder struct {
	ShadowOrder   *ShadowOrder
	Trades        []*tradeApp.Trade
	OriginalOrder *tradeApp.Order
}

func (e *ExpandOrder) Equal(_e *ExpandOrder) bool {
	if len(e.Trades) != len(_e.Trades) {
		return false
	}
	for i, trader := range e.Trades {
		if !trader.Equal(_e.Trades[i]) {
			return false
		}
	}

	return e.ShadowOrder.Equal(_e.ShadowOrder) &&
		e.OriginalOrder.Equal(_e.OriginalOrder)
}

func (e *ExpandOrder) Encode_Sign() ([]byte, error) {
	data := new(bytes.Buffer)

	// Shadow Order
	s, err := e.ShadowOrder.Encode_TransferBatching()
	if err != nil {
		return nil, err
	}
	hashedS := crypto.Keccak256Hash(s) // hash shadow order
	if err := binary.Write(data, binary.BigEndian, hashedS); err != nil {
		return nil, err
	}

	// Trades
	tData := new(bytes.Buffer)
	for _, trade := range e.Trades {
		t, err := trade.Encode_TransferBatching()
		if err != nil {
			return nil, err
		}
		if err := binary.Write(tData, binary.BigEndian, t); err != nil {
			return nil, err
		}
	}
	hashedTData := crypto.Keccak256Hash(tData.Bytes()) // hash trades
	if err := binary.Write(data, binary.BigEndian, hashedTData); err != nil {
		return nil, err
	}

	// Original Order
	o, err := e.OriginalOrder.Encode_TransferBatching()
	if err != nil {
		return nil, err
	}
	hashedO := crypto.Keccak256Hash(o)
	if err := binary.Write(data, binary.BigEndian, hashedO); err != nil {
		return nil, err
	}

	return data.Bytes(), err
}

func (e *ExpandOrder) IsValidOrder(owner common.Address) bool {
	amount := new(big.Int)
	for _, trade := range e.Trades {
		if e.ShadowOrder.From != e.OriginalOrder.OrderID &&
			trade.Owner.Cmp(owner) != 0 &&
			!trade.IsValidSignature() {
			return false
		}

		amount.Add(amount, trade.Amount)
	}
	_val := new(big.Int).Sub(e.OriginalOrder.Amount, amount)
	return _val.Cmp(e.ShadowOrder.Amount) == 0
}
