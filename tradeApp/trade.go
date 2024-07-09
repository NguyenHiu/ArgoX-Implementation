package tradeApp

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type Trade struct {
	TradeID   uuid.UUID
	BidOrder  uuid.UUID
	AskOrder  uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Owner     common.Address
	Signature []byte
}

func EndTrade(_ownerPrvKeyHex *ecdsa.PrivateKey) (*Trade, error) {
	trade := &Trade{
		TradeID:   EndID,
		BidOrder:  EndID,
		AskOrder:  EndID,
		Price:     new(big.Int),
		Amount:    new(big.Int),
		Owner:     crypto.PubkeyToAddress(_ownerPrvKeyHex.PublicKey),
		Signature: []byte{},
	}
	if err := trade.Sign(_ownerPrvKeyHex); err != nil {
		return nil, err
	}
	return trade, nil
}

func (t *Trade) Equal(_t *Trade) bool {
	return t.TradeID == _t.TradeID &&
		t.BidOrder == _t.BidOrder &&
		t.AskOrder == _t.AskOrder &&
		t.Amount.Cmp(_t.Amount) == 0 &&
		t.Price.Cmp(_t.Price) == 0 &&
		t.Owner.Cmp(_t.Owner) == 0 &&
		bytes.Equal(t.Signature, _t.Signature)
}

func (t *Trade) Clone() *Trade {
	clonedTrade := &Trade{
		TradeID:  t.TradeID,
		BidOrder: t.BidOrder,
		AskOrder: t.AskOrder,
		Owner:    t.Owner,
	}

	if t.Amount != nil {
		clonedTrade.Amount = new(big.Int).Set(t.Amount)
	}

	if t.Price != nil {
		clonedTrade.Price = new(big.Int).Set(t.Price)
	}

	if t.Signature != nil {
		clonedTrade.Signature = make([]byte, len(t.Signature))
		copy(clonedTrade.Signature, t.Signature)
	}

	return clonedTrade
}

func (t *Trade) Encode_Sign() ([]byte, error) {
	w := new(bytes.Buffer)

	// Trade ID
	tOrder, err := t.TradeID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, tOrder); err != nil {
		return nil, err
	}

	// BidOrder ID
	b, err := t.BidOrder.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, b); err != nil {
		return nil, err
	}

	// AskOrder ID
	s, err := t.AskOrder.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, s); err != nil {
		return nil, err
	}

	// Price
	if err := binary.Write(w, binary.BigEndian, PaddingToUint256(t.Price)); err != nil {
		return nil, err
	}

	// Amount
	if err := binary.Write(w, binary.BigEndian, PaddingToUint256(t.Amount)); err != nil {
		return nil, err
	}

	// Owner
	if err := binary.Write(w, binary.BigEndian, t.Owner.Bytes()); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (t *Trade) Sign(privateKey *ecdsa.PrivateKey) error {
	w, err := t.Encode_Sign()
	if err != nil {
		return err
	}

	// Sign
	hashedData := crypto.Keccak256Hash(w)
	sign, err := crypto.Sign(hashedData[:], privateKey)
	if err != nil {
		return err
	}

	t.Signature = sign

	return nil
}

func (t *Trade) Encode_TransferBatching() ([]byte, error) {
	w := new(bytes.Buffer)

	// Trade ID
	tOrder, err := t.TradeID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, tOrder); err != nil {
		return nil, err
	}

	// BidOrder ID
	b, err := t.BidOrder.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, b); err != nil {
		return nil, err
	}

	// AskOrder ID
	s, err := t.AskOrder.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, s); err != nil {
		return nil, err
	}

	// Price
	if err := binary.Write(w, binary.BigEndian, PaddingToUint256(t.Price)); err != nil {
		return nil, err
	}

	// Amount
	if err := binary.Write(w, binary.BigEndian, PaddingToUint256(t.Amount)); err != nil {
		return nil, err
	}

	// Owner
	if err := binary.Write(w, binary.BigEndian, t.Owner.Bytes()); err != nil {
		return nil, err
	}

	// Signature
	if err := binary.Write(w, binary.BigEndian, t.Signature); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (t *Trade) Decode_TransferBatching(data *bytes.Buffer) error {
	// Trade ID
	if err := binary.Read(data, binary.BigEndian, &t.TradeID); err != nil {
		return err
	}

	// BidOrder ID
	if err := binary.Read(data, binary.BigEndian, &t.BidOrder); err != nil {
		return err
	}

	// AskOrder ID
	if err := binary.Read(data, binary.BigEndian, &t.AskOrder); err != nil {
		return err
	}

	// Price
	_price := make([]byte, 32)
	if err := binary.Read(data, binary.BigEndian, &_price); err != nil {
		return err
	}
	t.Price = new(big.Int).SetBytes(_price)

	// Amount
	_amount := make([]byte, 32)
	if err := binary.Read(data, binary.BigEndian, &_amount); err != nil {
		return err
	}
	t.Amount = new(big.Int).SetBytes(_amount)

	// Owner
	_owner := make([]byte, 20)
	if err := binary.Read(data, binary.BigEndian, &_owner); err != nil {
		return err
	}
	t.Owner = common.Address(_owner)

	// Signature
	_signature := make([]byte, 65)
	if err := binary.Read(data, binary.BigEndian, &_signature); err != nil {
		return err
	}
	t.Signature = _signature

	return nil
}

func (t *Trade) IsValidSignature() bool {
	encodedTrade, err := t.Encode_Sign()
	if err != nil {
		return false
	}

	hashedData := crypto.Keccak256Hash(encodedTrade)
	pubkey, err := crypto.SigToPub(hashedData.Bytes(), t.Signature)
	if err != nil {
		return false
	}

	if crypto.PubkeyToAddress(*pubkey).Cmp(t.Owner) != 0 {
		return false
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(pubkey), hashedData.Bytes(), t.Signature[:64])
}

func (t *Trade) Encode_TransferLightning() ([]byte, error) {
	return t.Encode_TransferBatching()
}

func (t *Trade) Decode_TransferLightning(data *bytes.Buffer) error {
	return t.Decode_TransferBatching(data)
}
