package matcher

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
	BidOrder  uuid.UUID
	AskOrder  uuid.UUID
	Amount    *big.Int
	Owner     common.Address
	Signature []byte
}

func (t *Trade) Equal(_t *Trade) bool {
	return t.BidOrder == _t.BidOrder &&
		t.AskOrder == _t.AskOrder &&
		t.Amount.Cmp(_t.Amount) == 0 &&
		t.Owner.Cmp(_t.Owner) == 0 &&
		bytes.Compare(t.Signature, _t.Signature) == 0
}

func (t *Trade) Encode_Sign() ([]byte, error) {
	w := new(bytes.Buffer)

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
	// BidOrder ID
	if err := binary.Read(data, binary.BigEndian, &t.BidOrder); err != nil {
		return err
	}

	// AskOrder ID
	if err := binary.Read(data, binary.BigEndian, &t.AskOrder); err != nil {
		return err
	}

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
