package matcher

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type Batch struct {
	BatchID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Orders    []*MatcherOrder
	Owner     common.Address
	Signature []byte
}

func newBatch(_price, _amount *big.Int, _side bool, _orders []*MatcherOrder, _owner common.Address) *Batch {
	id, _ := uuid.NewRandom()
	return &Batch{
		BatchID: id,
		Price:   _price,
		Amount:  _amount,
		Side:    _side,
		Orders:  _orders,
		Owner:   _owner,
	}
}

func (b *Batch) Sign(privateKey string) {
	// Get encoded data of all orders
	ordersData := new(bytes.Buffer)
	for _, order := range b.Orders {
		binary.Write(ordersData, binary.BigEndian, order.Data.EncodePackedOrder())
	}

	// Encode packed the batch
	data := new(bytes.Buffer)
	binary.Write(data, binary.BigEndian, b.BatchID)
	binary.Write(data, binary.BigEndian, b.Price)
	binary.Write(data, binary.BigEndian, b.Amount)
	binary.Write(data, binary.BigEndian, b.Side)
	binary.Write(data, binary.BigEndian, ordersData.Bytes())
	binary.Write(data, binary.BigEndian, b.Owner)

	// Hash the encode packed data
	hasheddata := crypto.Keccak256Hash(data.Bytes())

	_prvkey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the batch
	sig, err := crypto.Sign(hasheddata.Bytes(), _prvkey)
	if err != nil {
		log.Fatal(err)
	}

	b.Signature = sig
}

// I know this function isn't ideal and has too much duplicated code.
// However, the current priority is to get the protocol running.
func (m *Matcher) batching() ([]*Batch, []*Batch) {

	// bid batches
	bidBatches := []*Batch{}
	if len(m.BidOrders) > 0 {
		ord := m.BidOrders[0]
		m.BidOrders = m.BidOrders[1:]
		batch := newBatch(ord.Data.Price, ord.Data.Amount, ord.Data.Side, []*MatcherOrder{ord}, m.Address)
		cnt := 1
		for cnt <= constants.NO_BATCHES_EACH_TIME && len(m.BidOrders) > 0 {
			ord = m.BidOrders[0]
			if ord.Data.Price.Cmp(batch.Price) == 0 {
				m.BidOrders = m.BidOrders[1:]
				batch.Amount = new(big.Int).Add(batch.Amount, ord.Data.Amount)
				batch.Orders = append(batch.Orders, ord)
			} else {
				cnt++
				bidBatches = append(bidBatches, batch)
				batch = newBatch(ord.Data.Price, ord.Data.Amount, ord.Data.Side, []*MatcherOrder{ord}, m.Address)
			}
		}
		if cnt < constants.NO_BATCHES_EACH_TIME {
			for len(m.BidOrders) > 0 {
				ord = m.BidOrders[0]
				if ord.Data.Price.Cmp(batch.Price) == 0 {
					m.BidOrders = m.BidOrders[1:]
					batch.Amount = new(big.Int).Add(batch.Amount, ord.Data.Amount)
					batch.Orders = append(batch.Orders, ord)
				} else {
					break
				}
			}
			bidBatches = append(bidBatches, batch)
		}
	}

	// ask batches
	askBatches := []*Batch{}
	if len(m.AskOrders) > 0 {
		ord := m.AskOrders[0]
		m.AskOrders = m.AskOrders[1:]
		batch := newBatch(ord.Data.Price, ord.Data.Amount, ord.Data.Side, []*MatcherOrder{ord}, m.Address)
		cnt := 1
		for cnt <= constants.NO_BATCHES_EACH_TIME && len(m.AskOrders) > 0 {
			ord = m.AskOrders[0]
			if ord.Data.Price.Cmp(batch.Price) == 0 {
				m.AskOrders = m.AskOrders[1:]
				batch.Amount = new(big.Int).Add(batch.Amount, ord.Data.Amount)
				batch.Orders = append(batch.Orders, ord)
			} else {
				cnt++
				askBatches = append(askBatches, batch)
				batch = newBatch(ord.Data.Price, ord.Data.Amount, ord.Data.Side, []*MatcherOrder{ord}, m.Address)
			}
		}
		if cnt < constants.NO_BATCHES_EACH_TIME {
			for len(m.AskOrders) > 0 {
				ord = m.AskOrders[0]
				if ord.Data.Price.Cmp(batch.Price) == 0 {
					m.AskOrders = m.AskOrders[1:]
					batch.Amount = new(big.Int).Add(batch.Amount, ord.Data.Amount)
					batch.Orders = append(batch.Orders, ord)
				} else {
					break
				}
			}
			askBatches = append(askBatches, batch)
		}
	}

	return bidBatches, askBatches
}
