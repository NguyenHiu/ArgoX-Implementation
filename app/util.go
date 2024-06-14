package app

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
)

type Order struct {
	OrderID       uuid.UUID
	Price         *big.Int
	Amount        *big.Int
	Side          bool
	Owner         *wallet.Address
	OwnerSignture []byte
	Status        string
	MatchedAmount *big.Int
}

type OrderUpdatedInfo struct {
	Status        string
	MatchedAmount *big.Int
}

// The `status` parameter should be "P" at the init phase,
// but allowing the `status` parameter to be passed is for testing purposes.
func NewOrder(price, amount *big.Int, side bool, owner *wallet.Address, status string) Order {
	orderId, _ := uuid.NewRandom()
	return Order{
		OrderID:       orderId,
		Price:         price,
		Amount:        amount,
		Side:          side,
		Owner:         owner,
		OwnerSignture: []byte{},
		Status:        status, // Replace later
		MatchedAmount: &big.Int{},
	}
}

func (o *Order) Sign(prvkey ecdsa.PrivateKey) error {
	pub, _ := prvkey.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)
	if o.Owner.Cmp(wallet.AsWalletAddr(addr)) != 0 {
		return fmt.Errorf("private key does not match with the order's owner")
	}
	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("invalid uuid")
	}
	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, orderID)
	binary.Write(data, binary.BigEndian, PaddingToUint256(o.Price))
	binary.Write(data, binary.BigEndian, PaddingToUint256(o.Amount))
	binary.Write(data, binary.BigEndian, o.Side)
	binary.Write(data, binary.BigEndian, o.Owner.Bytes())

	hashedData := crypto.Keccak256Hash(data.Bytes())

	sig, err := crypto.Sign(hashedData.Bytes(), &prvkey)
	if err != nil {
		return fmt.Errorf("can not sign the order, err: %v", err)
	}
	o.OwnerSignture = sig

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

	pub, err := crypto.SigToPub(hashedData.Bytes(), o.OwnerSignture)
	if err != nil {
		fmt.Printf("Cannot recover public key from signature, error: %v\n", err)
		return false
	}
	_owner := wallet.AsWalletAddr(crypto.PubkeyToAddress(*pub))
	if _owner.Cmp(o.Owner) != 0 {
		fmt.Println("Provided public key does not match with the order's owner")
		return false
	}
	pubBytes := crypto.FromECDSAPub(pub)
	return crypto.VerifySignature(pubBytes, hashedData.Bytes(), o.OwnerSignture[:64])

}

func (o *Order) Equal(_o *Order) bool {
	return (o.OrderID == _o.OrderID &&
		o.Price.Cmp(_o.Price) == 0 &&
		o.Amount.Cmp(_o.Amount) == 0 &&
		o.Side == _o.Side &&
		o.Owner.Cmp(_o.Owner) == 0 &&
		o.Status == _o.Status &&
		o.MatchedAmount == _o.MatchedAmount)
}

// Encode Order
// Price > Amount > Side > Owner > Status
func (o *Order) EncodeOrder() []byte {
	buf := new(bytes.Buffer)

	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		fmt.Println("invalid uuid")
	}
	err = binary.Write(buf, binary.BigEndian, orderID)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Price))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Amount))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Side)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Owner.Bytes())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.OwnerSignture)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, []byte(o.Status))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.MatchedAmount))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func (o *Order) EncodePackedOrder() []byte {
	buf := new(bytes.Buffer)

	orderID, err := o.OrderID.MarshalBinary()
	if err != nil {
		fmt.Println("invalid uuid")
	}
	err = binary.Write(buf, binary.BigEndian, orderID)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Price))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, PaddingToUint256(o.Amount))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Side)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	err = binary.Write(buf, binary.BigEndian, o.Owner.Bytes())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

// Decode Order
// Follow the parameter orders when encoding
func DecodeOrder(data []byte) (*Order, error) {
	order := Order{}
	buf := bytes.NewBuffer(data)

	orderIDTemp := make([]byte, 16)
	err := binary.Read(buf, binary.BigEndian, &orderIDTemp)
	if err != nil {
		return nil, err
	}
	err = order.OrderID.UnmarshalBinary(orderIDTemp)
	if err != nil {
		return nil, err
	}

	price := make([]byte, 32)
	err = binary.Read(buf, binary.BigEndian, &price)
	if err != nil {
		return nil, err
	}
	order.Price = new(big.Int).SetBytes(price)

	amount := make([]byte, 32)
	err = binary.Read(buf, binary.BigEndian, &amount)
	if err != nil {
		return nil, err
	}
	order.Amount = new(big.Int).SetBytes(amount)

	err = binary.Read(buf, binary.BigEndian, &order.Side)
	if err != nil {
		return nil, err
	}

	order.Owner = &wallet.Address{}
	err = binary.Read(buf, binary.BigEndian, &order.Owner)
	if err != nil {
		return nil, err
	}

	ownerSign := make([]byte, 65)
	err = binary.Read(buf, binary.BigEndian, &ownerSign)
	if err != nil {
		return nil, err
	}
	order.OwnerSignture = ownerSign

	status_temp := make([]byte, 1)
	err = binary.Read(buf, binary.BigEndian, &status_temp)
	if err != nil {
		return nil, err
	}
	order.Status = string(status_temp)

	matchedAmount := make([]byte, 32)
	err = binary.Read(buf, binary.BigEndian, &matchedAmount)
	if err != nil {
		return nil, err
	}
	order.MatchedAmount = new(big.Int).SetBytes(matchedAmount)

	return &order, nil
}

func (d *VerifyAppData) CheckFinal() bool {
	l := len(d.Orders)
	return l != 0 && d.Orders[l-1].Status == "F"
}

func computeFinalBalances(orders []*Order, bals channel.Balances) channel.Balances {
	matcherReceivedETH := &big.Int{}
	matcherReceivedGAV := &big.Int{}

	for i := 0; i < len(orders); i++ {
		if orders[i].Status != "F" && orders[i].Status == "M" {
			if orders[i].Side == constants.BID {
				matcherReceivedETH = new(big.Int).Add(matcherReceivedETH, orders[i].Price)
				matcherReceivedGAV = new(big.Int).Sub(matcherReceivedGAV, orders[i].Amount)
			} else {
				matcherReceivedETH = new(big.Int).Sub(matcherReceivedETH, orders[i].Price)
				matcherReceivedGAV = new(big.Int).Add(matcherReceivedGAV, orders[i].Amount)
			}
		}
	}

	finalBals := bals.Clone()
	finalBals[constants.ETH][constants.MATCHER] = new(big.Int).Add(bals[constants.ETH][constants.MATCHER], matcherReceivedETH)
	finalBals[constants.ETH][constants.TRADER] = new(big.Int).Sub(bals[constants.ETH][constants.TRADER], matcherReceivedETH)
	finalBals[constants.GVN][constants.MATCHER] = new(big.Int).Add(bals[constants.GVN][constants.MATCHER], matcherReceivedGAV)
	finalBals[constants.GVN][constants.TRADER] = new(big.Int).Sub(bals[constants.GVN][constants.TRADER], matcherReceivedGAV)

	return finalBals
}

func PaddingToUint256(num *big.Int) []byte {
	return append(make([]byte, 32-len(num.Bytes())), num.Bytes()...)
}
