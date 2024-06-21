package matcher

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// 16 + 32 + 32 + 1 + 1 + x*166 + 20 + 65
type Batch struct {
	BatchID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Orders    []*MatcherOrder
	Owner     common.Address
	Signature []byte
}

// Encode batch to the format of super matcher batch
func (b *Batch) Encode_TranferBatching() ([]byte, error) {
	data := new(bytes.Buffer)
	id, err := b.BatchID.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}
	if err := binary.Write(data, binary.BigEndian, id); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Price)); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Amount)); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Side); err != nil {
		return []byte{}, err
	}

	ordersbyte := new(bytes.Buffer)
	for _, order := range b.Orders {
		encodedOrder, err := order.Data.Encode_TransferBatching()
		if err != nil {
			return []byte{}, err
		}

		// TODO: clear
		if len(encodedOrder) != 166 {
			_logger.Debug("len(encodedOrder6 is not equal to 166 bytes)")
		}

		if err := binary.Write(ordersbyte, binary.BigEndian, encodedOrder); err != nil {
			return []byte{}, err
		}
	}

	if err := binary.Write(data, binary.BigEndian, uint8(len(b.Orders))); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, ordersbyte.Bytes()); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Owner.Bytes()); err != nil {
		return []byte{}, err
	}

	if err := binary.Write(data, binary.BigEndian, b.Signature); err != nil {
		return []byte{}, err
	}

	return data.Bytes(), nil
}

// func newBatch(_price, _amount *big.Int, _side bool, _orders []*MatcherOrder, _owner common.Address) *Batch {
// 	id, _ := uuid.NewRandom()
// 	return &Batch{
// 		BatchID: id,
// 		Price:   _price,
// 		Amount:  _amount,
// 		Side:    _side,
// 		Orders:  _orders,
// 		Owner:   _owner,
// 	}
// }

func (b *Batch) Sign(_prvkey *ecdsa.PrivateKey) error {
	// Get encoded data of all orders
	ordersData := new(bytes.Buffer)
	for _, order := range b.Orders {
		data, err := order.Data.Encode_TransferBatching()
		if err != nil {
			return fmt.Errorf("invalid order data in Sign() func, err: %v", err)
		}
		binary.Write(ordersData, binary.BigEndian, data)
	}

	// Encode packed the batch
	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, b.BatchID)
	binary.Write(data, binary.BigEndian, PaddingToUint256(b.Price))
	binary.Write(data, binary.BigEndian, PaddingToUint256(b.Amount))
	binary.Write(data, binary.BigEndian, b.Side)
	binary.Write(data, binary.BigEndian, uint8(len(b.Orders)))
	binary.Write(data, binary.BigEndian, ordersData.Bytes())
	binary.Write(data, binary.BigEndian, b.Owner)

	// Hash the encode packed data
	hasheddata := crypto.Keccak256Hash(data.Bytes())

	// _prvkey, err := crypto.HexToECDSA(privateKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Sign the batch
	sig, err := crypto.Sign(hasheddata.Bytes(), _prvkey)
	if err != nil {
		log.Fatal(err)
	}

	b.Signature = sig

	return nil
}

func (m *Matcher) batching() []*Batch {
	batches := []*Batch{}

	if len(m.BidOrders) > 0 {
		ord := m.BidOrders[0]
		m.BidOrders = m.BidOrders[1:]
		id, _ := uuid.NewRandom()
		batch := &Batch{
			BatchID: id,
			Price:   ord.Data.Price,
			Amount:  ord.Data.Amount,
			Side:    ord.Data.Side,
			Orders:  []*MatcherOrder{ord},
			Owner:   m.Address,
		}
		batch.Sign(m.PrivateKey)
		batches = append(batches, batch)
	}

	if len(m.AskOrders) > 0 {
		ord := m.AskOrders[0]
		m.AskOrders = m.AskOrders[1:]
		id, _ := uuid.NewRandom()
		batch := &Batch{
			BatchID: id,
			Price:   ord.Data.Price,
			Amount:  ord.Data.Amount,
			Side:    ord.Data.Side,
			Orders:  []*MatcherOrder{ord},
			Owner:   m.Address,
		}
		batch.Sign(m.PrivateKey)
		batches = append(batches, batch)
	}

	return batches
}

func PaddingToUint256(num *big.Int) []byte {
	return append(make([]byte, 32-len(num.Bytes())), num.Bytes()...)
}
