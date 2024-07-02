package matcher

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/google/uuid"
)

type ShadowOrder struct {
	Price  *big.Int
	Amount *big.Int
	Side   bool
	From   uuid.UUID
}

func (o *ShadowOrder) Encode_TransferBatching() ([]byte, error) {
	data := new(bytes.Buffer)

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(o.Price)); err != nil {
		return nil, err
	}
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(o.Amount)); err != nil {
		return nil, err
	}
	if err := binary.Write(data, binary.BigEndian, o.Side); err != nil {
		return nil, err
	}
	id, err := o.From.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(data, binary.BigEndian, id); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (o *ShadowOrder) Decode_TransferBatching(data *bytes.Buffer) error {
	// Price
	_price := make([]byte, 32)
	if err := binary.Read(data, binary.BigEndian, &_price); err != nil {
		return err
	}
	o.Price = new(big.Int).SetBytes(_price)

	// Amount
	_amount := make([]byte, 32)
	if err := binary.Read(data, binary.BigEndian, &_amount); err != nil {
		return err
	}
	o.Amount = new(big.Int).SetBytes(_amount)

	// Side
	if err := binary.Read(data, binary.BigEndian, &o.Side); err != nil {
		return err
	}

	// From
	if err := binary.Read(data, binary.BigEndian, &o.From); err != nil {
		return err
	}

	return nil
}

func (o *ShadowOrder) Equal(_o *ShadowOrder) bool {
	return o.Price.Cmp(_o.Price) == 0 &&
		o.Amount.Cmp(_o.Amount) == 0 &&
		o.Side == _o.Side &&
		o.From == _o.From
}
