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

func (m *ShadowOrder) Clone() *ShadowOrder {
	return &ShadowOrder{
		Price:  new(big.Int).Set(m.Price),
		Amount: new(big.Int).Set(m.Amount),
		Side:   m.Side,
		From:   m.From,
	}
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

func (o *ShadowOrder) Equal(_o *ShadowOrder) bool {
	return o.Price.Cmp(_o.Price) == 0 &&
		o.Amount.Cmp(_o.Amount) == 0 &&
		o.Side == _o.Side &&
		o.From == _o.From
}
