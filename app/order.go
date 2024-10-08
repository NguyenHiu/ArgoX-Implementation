package app

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
)

var EndID = uuid.UUID{}

type Order struct {
	OrderID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Owner     *wallet.Address
	Signature []byte
}

func (o *Order) Equal(_o *Order) bool {
	return o.OrderID == _o.OrderID &&
		o.Price.Cmp(_o.Price) == 0 &&
		o.Amount.Cmp(_o.Amount) == 0 &&
		o.Side == _o.Side &&
		o.Owner.Cmp(_o.Owner) == 0 &&
		bytes.Equal(o.Signature, _o.Signature)
}

func (o *Order) Clone() *Order {
	newOrder := &Order{
		OrderID:   o.OrderID,
		Side:      o.Side,
		Signature: make([]byte, len(o.Signature)),
	}

	copy(newOrder.Signature, o.Signature)

	if o.Price != nil {
		newOrder.Price = new(big.Int).Set(o.Price)
	}
	if o.Amount != nil {
		newOrder.Amount = new(big.Int).Set(o.Amount)
	}

	if o.Owner != nil {
		newOwner := *o.Owner
		newOrder.Owner = &newOwner
	}

	return newOrder
}

func NewOrder(price, amount *big.Int, side bool, owner *wallet.Address) *Order {
	orderId, _ := uuid.NewRandom()
	return &Order{
		OrderID:   orderId,
		Price:     price,
		Amount:    amount,
		Side:      side,
		Owner:     owner,
		Signature: []byte{},
	}
}

func EndOrder(_ownerPrvKeyHex string) (*Order, error) {
	prvkey, err := crypto.HexToECDSA(_ownerPrvKeyHex)
	if err != nil {
		return nil, err
	}

	newOrder := &Order{
		OrderID:   uuid.UUID{},
		Price:     new(big.Int),
		Amount:    new(big.Int),
		Side:      true,
		Owner:     wallet.AsWalletAddr(crypto.PubkeyToAddress(prvkey.PublicKey)),
		Signature: []byte{},
	}
	if err := newOrder.Sign(_ownerPrvKeyHex); err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (o *Order) Sign(privateKey string) error {
	prvKey, _ := crypto.HexToECDSA(privateKey)
	addr := crypto.PubkeyToAddress(prvKey.PublicKey)
	if o.Owner.Cmp(wallet.AsWalletAddr(addr)) != 0 {
		return fmt.Errorf("private key does not match with the order's owner")
	}

	data, err := o.Encode_Sign()
	if err != nil {
		return err
	}

	hashedData := crypto.Keccak256Hash(data)

	sig, err := crypto.Sign(hashedData.Bytes(), prvKey)
	if err != nil {
		return fmt.Errorf("can not sign the order, err: %v", err)
	}
	o.Signature = sig

	return nil
}

func (o *Order) IsValidSignature() bool {
	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		return false
	}
	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, orderID)
	binary.Write(data, binary.BigEndian, PaddingToUint256(o.Price))
	binary.Write(data, binary.BigEndian, PaddingToUint256(o.Amount))
	binary.Write(data, binary.BigEndian, o.Side)
	binary.Write(data, binary.BigEndian, o.Owner.Bytes())
	hashedData := crypto.Keccak256Hash(data.Bytes())

	pub, err := crypto.SigToPub(hashedData.Bytes(), o.Signature)
	if err != nil {
		_logger.Debug("Cannot recover public key from signature, error: %v\n", err)
		return false
	}
	_owner := wallet.AsWalletAddr(crypto.PubkeyToAddress(*pub))
	if _owner.Cmp(o.Owner) != 0 {
		_logger.Debug("Provided public key does not match with the order's owner\n")
		return false
	}
	pubBytes := crypto.FromECDSAPub(pub)
	return crypto.VerifySignature(pubBytes, hashedData.Bytes(), o.Signature[:64])
}

// Used in Lightning
func (o *Order) Encode_TransferLightning() []byte {
	buf := new(bytes.Buffer)

	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		_logger.Debug("invalid uuid\n")
	}
	err = binary.Write(buf, binary.BigEndian, orderID)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Price))
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Amount))
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Side)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Owner.Bytes())
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Signature)
	if err != nil {
		_logger.Debug("binary.Write failed: %v\n", err)
	}

	return buf.Bytes()
}

// Decode Order
// Follow the parameter orders when encoding
func (o *Order) Decode_TransferLightning(buf *bytes.Buffer) error {

	orderIDTemp := make([]byte, 16)
	err := binary.Read(buf, binary.BigEndian, &orderIDTemp)
	if err != nil {
		return err
	}
	err = o.OrderID.UnmarshalBinary(orderIDTemp)
	if err != nil {
		return err
	}

	price := make([]byte, 32)
	err = binary.Read(buf, binary.BigEndian, &price)
	if err != nil {
		return err
	}
	o.Price = new(big.Int).SetBytes(price)

	amount := make([]byte, 32)
	err = binary.Read(buf, binary.BigEndian, &amount)
	if err != nil {
		return err
	}
	o.Amount = new(big.Int).SetBytes(amount)

	err = binary.Read(buf, binary.BigEndian, &o.Side)
	if err != nil {
		return err
	}

	owner := make([]byte, 20)
	err = binary.Read(buf, binary.BigEndian, &owner)
	if err != nil {
		return err
	}
	o.Owner = (*wallet.Address)(owner)

	ownerSign := make([]byte, 65)
	err = binary.Read(buf, binary.BigEndian, &ownerSign)
	if err != nil {
		return err
	}
	o.Signature = ownerSign

	return nil
}

// Used in Batching
func (o *Order) Encode_TransferBatching() ([]byte, error) {
	buf := new(bytes.Buffer)

	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		return []byte{}, fmt.Errorf("invalid uuid")
	}
	err = binary.Write(buf, binary.BigEndian, orderID)
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed: %v", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Price))
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed: %v", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Amount))
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed: %v", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Side)
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed: %v", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Owner.Bytes())
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed: %v", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Signature)
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed: %v", err)
	}

	return buf.Bytes(), nil
}

func (o *Order) Decode_TransferBatching(data *bytes.Buffer) error {

	orderIDTemp := make([]byte, 16)
	err := binary.Read(data, binary.BigEndian, &orderIDTemp)
	if err != nil {
		return err
	}
	err = o.OrderID.UnmarshalBinary(orderIDTemp)
	if err != nil {
		return err
	}

	price := make([]byte, 32)
	err = binary.Read(data, binary.BigEndian, &price)
	if err != nil {
		return err
	}
	o.Price = new(big.Int).SetBytes(price)

	amount := make([]byte, 32)
	err = binary.Read(data, binary.BigEndian, &amount)
	if err != nil {
		return err
	}
	o.Amount = new(big.Int).SetBytes(amount)

	err = binary.Read(data, binary.BigEndian, &o.Side)
	if err != nil {
		return err
	}

	owner := make([]byte, 20)
	err = binary.Read(data, binary.BigEndian, &owner)
	if err != nil {
		return err
	}
	o.Owner = (*wallet.Address)(owner)

	ownerSign := make([]byte, 65)
	err = binary.Read(data, binary.BigEndian, &ownerSign)
	if err != nil {
		return err
	}
	o.Signature = ownerSign

	return nil
}

// Used in Smart Contract
func (o *Order) Encode_Sign() ([]byte, error) {
	buf := new(bytes.Buffer)

	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		return []byte{}, fmt.Errorf("invalid uuid")
	}
	err = binary.Write(buf, binary.BigEndian, orderID)
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed:%v", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Price))
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed:%v", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Amount))
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed:%v", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Side)
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed:%v", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Owner.Bytes())
	if err != nil {
		return []byte{}, fmt.Errorf("binary.Write failed:%v", err)
	}

	return buf.Bytes(), nil
}
