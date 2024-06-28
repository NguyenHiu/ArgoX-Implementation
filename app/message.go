package app

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
)

type Message struct {
	MessageID     uuid.UUID
	OrderID       uuid.UUID
	MatchedAmount *big.Int
	Status        uint8
	Owner         *wallet.Address
	Signature     []byte
}

func (m *Message) Clone() *Message {
	// Deep copy for MatchedAmount
	matchedAmountCopy := new(big.Int).Set(m.MatchedAmount)

	// Assuming wallet.Address is immutable, we just copy the pointer
	// If it's not, you would need to implement a deep copy method for it

	// Deep copy for Signature
	signatureCopy := make([]byte, len(m.Signature))
	copy(signatureCopy, m.Signature)

	return &Message{
		MessageID:     m.MessageID,
		OrderID:       m.OrderID,
		MatchedAmount: matchedAmountCopy,
		Status:        m.Status,
		Owner:         m.Owner, // Directly copy the pointer if immutable
		Signature:     signatureCopy,
	}
}

func NewMsg(_orderID uuid.UUID, _matchedAmount *big.Int, _status uint8, _owner common.Address) Message {
	messageID, _ := uuid.NewRandom()
	return Message{
		MessageID:     messageID,
		OrderID:       _orderID,
		MatchedAmount: _matchedAmount,
		Status:        _status,
		Owner:         wallet.AsWalletAddr(_owner),
		Signature:     []byte{},
	}
}

func (m *Message) Sign(prvkey ecdsa.PrivateKey) error {
	pub, _ := prvkey.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)
	if m.Owner.Cmp(wallet.AsWalletAddr(addr)) != 0 {
		return fmt.Errorf("private key does not match with the order's owner")
	}

	messageID, err := m.MessageID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("invalid uuid")
	}
	orderID, err := m.OrderID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("invalid uuid")
	}
	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, messageID)
	binary.Write(data, binary.BigEndian, orderID)
	binary.Write(data, binary.BigEndian, PaddingToUint256(m.MatchedAmount))
	binary.Write(data, binary.BigEndian, m.Status)
	binary.Write(data, binary.BigEndian, m.Owner.Bytes())

	hashedData := crypto.Keccak256Hash(data.Bytes())

	sig, err := crypto.Sign(hashedData.Bytes(), &prvkey)
	if err != nil {
		return fmt.Errorf("can not sign the order, err: %v", err)
	}
	m.Signature = sig

	return nil
}

func (m *Message) IsValidSignature() bool {
	messageID, err := m.OrderID.MarshalBinary()
	if err != nil {
		return false
	}
	orderID, err := m.OrderID.MarshalBinary()
	if err != nil {
		return false
	}
	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, messageID)
	binary.Write(data, binary.BigEndian, orderID)
	binary.Write(data, binary.BigEndian, PaddingToUint256(m.MatchedAmount))
	binary.Write(data, binary.BigEndian, m.Status)
	binary.Write(data, binary.BigEndian, m.Owner.Bytes())
	hashedData := crypto.Keccak256Hash(data.Bytes())

	pub, err := crypto.SigToPub(hashedData.Bytes(), m.Signature)
	if err != nil {
		_logger.Debug("Cannot recover public key from signature, error: %v\n", err)
		return false
	}
	_owner := wallet.AsWalletAddr(crypto.PubkeyToAddress(*pub))
	if _owner.Cmp(m.Owner) != 0 {
		_logger.Debug("Provided public key does not match with the order's owner\n")
		return false
	}
	pubBytes := crypto.FromECDSAPub(pub)
	return crypto.VerifySignature(pubBytes, hashedData.Bytes(), m.Signature[:64])
}

func (m *Message) Equal(_m *Message) bool {
	return (m.MessageID == _m.MessageID &&
		m.OrderID == _m.OrderID &&
		m.MatchedAmount.Cmp(_m.MatchedAmount) == 0 &&
		m.Status == _m.Status &&
		m.Owner.Cmp(_m.Owner) == 0 &&
		bytes.Equal(m.Signature, _m.Signature))
}

// Used in Lightning
func (m *Message) Encode_TransferLightning() []byte {
	buf := new(bytes.Buffer)

	messageID, err := m.MessageID.MarshalBinary()
	if err != nil {
		_logger.Debug("invalid uuid\n")
	}
	err = binary.Write(buf, binary.BigEndian, messageID)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	orderID, err := m.OrderID.MarshalBinary()
	if err != nil {
		_logger.Debug("invalid uuid\n")
	}
	err = binary.Write(buf, binary.BigEndian, orderID)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(m.MatchedAmount))
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, m.Status)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, m.Owner.Bytes())
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, m.Signature)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}
	return buf.Bytes()
}

// Decode Message
// Follow the parameter orders when encoding
func Message_Decode_TransferLightning(data []byte) (*Message, error) {
	message := Message{}
	buf := bytes.NewBuffer(data)

	messageIDTemp := make([]byte, 16)
	err := binary.Read(buf, binary.BigEndian, &messageIDTemp)
	if err != nil {
		return nil, err
	}
	err = message.MessageID.UnmarshalBinary(messageIDTemp)
	if err != nil {
		return nil, err
	}

	orderIDTemp := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &orderIDTemp)
	if err != nil {
		return nil, err
	}
	err = message.OrderID.UnmarshalBinary(orderIDTemp)
	if err != nil {
		return nil, err
	}

	matchedAmount := make([]byte, 32)
	err = binary.Read(buf, binary.BigEndian, &matchedAmount)
	if err != nil {
		return nil, err
	}
	message.MatchedAmount = new(big.Int).SetBytes(matchedAmount)

	err = binary.Read(buf, binary.BigEndian, &message.Status)
	if err != nil {
		return nil, err
	}

	owner := make([]byte, 20)
	err = binary.Read(buf, binary.BigEndian, &owner)
	if err != nil {
		return nil, err
	}
	message.Owner = (*wallet.Address)(owner)

	ownerSign := make([]byte, 65)
	err = binary.Read(buf, binary.BigEndian, &ownerSign)
	if err != nil {
		return nil, err
	}
	message.Signature = ownerSign

	return &message, nil
}
