package matcher

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/tradeApp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type MatcherOrder struct {
	Data  *ShadowOrder
	Owner uuid.UUID
}

func (m *MatcherOrder) Clone() *MatcherOrder {
	return &MatcherOrder{
		Data:  m.Data.Clone(),
		Owner: m.Owner,
	}
}

type Batch struct {
	BatchID   uuid.UUID
	Price     *big.Int
	Amount    *big.Int
	Side      bool
	Orders    []*ExpandOrder
	Owner     common.Address
	Signature []byte
}

func (m *Matcher) NewBatch(_price, _amount *big.Int, _side bool, _orders []*MatcherOrder) *Batch {
	id, _ := uuid.NewRandom()
	orders := []*ExpandOrder{}
	for _, order := range _orders {
		var trades []*tradeApp.Trade
		if order.Data.Side == constants.BID {
			trades = m.mappingBidtoTrade[order.Data.From]
		} else {
			trades = m.mappingAskToTrade[order.Data.From]
		}
		orders = append(orders, &ExpandOrder{
			ShadowOrder:   order.Data,
			Trades:        trades,
			OriginalOrder: m.Orders[order.Data.From],
		})
	}

	batch := &Batch{
		BatchID:   id,
		Price:     _price,
		Amount:    _amount,
		Side:      _side,
		Orders:    orders,
		Owner:     m.Address,
		Signature: []byte{},
	}

	return batch
}

// Encode batch to the format of super matcher batch
func (b *Batch) Encode_TransferBatching(m *Matcher) ([]byte, error) {
	data := new(bytes.Buffer)

	// Batch ID
	id, err := b.BatchID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(data, binary.BigEndian, id); err != nil {
		return nil, err
	}

	// Price
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Price)); err != nil {
		return nil, err
	}

	// Amount
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Amount)); err != nil {
		return nil, err
	}

	// Side
	if err := binary.Write(data, binary.BigEndian, b.Side); err != nil {
		return nil, err
	}

	// Number of orders
	if err := binary.Write(data, binary.BigEndian, uint8(len(b.Orders))); err != nil {
		return nil, err
	}

	// Current Order + No. Executed Trades + Executed Trades + Original Order
	ordersbyte := new(bytes.Buffer)
	for _, order := range b.Orders {
		// Shadow order
		order_, err := order.ShadowOrder.Encode_TransferBatching()
		if err != nil {
			return nil, err
		}
		if err := binary.Write(ordersbyte, binary.BigEndian, order_); err != nil {
			return nil, err
		}

		// No. Executed Trades
		var executedTrades []*tradeApp.Trade
		if order.ShadowOrder.Side == constants.BID {
			executedTrades = m.mappingBidtoTrade[order.ShadowOrder.From]
		} else {
			executedTrades = m.mappingAskToTrade[order.ShadowOrder.From]
		}
		if err := binary.Write(ordersbyte, binary.BigEndian, uint8(len(executedTrades))); err != nil {
			return nil, err
		}

		// Executed Trades
		for _, trade := range executedTrades {
			d, err := trade.Encode_TransferBatching()
			if err != nil {
				return nil, err
			}
			if err := binary.Write(ordersbyte, binary.BigEndian, d); err != nil {
				return nil, err
			}
		}

		// Original Order
		d, err := m.Orders[order.ShadowOrder.From].Encode_TransferBatching()
		if err != nil {
			return nil, err
		}
		if err := binary.Write(ordersbyte, binary.BigEndian, d); err != nil {
			return nil, err
		}
	}
	if err := binary.Write(data, binary.BigEndian, ordersbyte.Bytes()); err != nil {
		return nil, err
	}

	// Owner
	if err := binary.Write(data, binary.BigEndian, b.Owner.Bytes()); err != nil {
		return nil, err
	}

	// Signature
	if err := binary.Write(data, binary.BigEndian, b.Signature); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (b *Batch) Encode_Sign() ([]byte, error) {
	// Get encoded data of all orders
	ordersData := new(bytes.Buffer)
	for _, order := range b.Orders {
		data, err := order.Encode_Sign()
		if err != nil {
			return nil, fmt.Errorf("invalid order data in Sign() func, err: %v", err)
		}
		binary.Write(ordersData, binary.BigEndian, data)
	}

	// Encode packed the batch
	data := new(bytes.Buffer)

	// Batch ID
	if err := binary.Write(data, binary.BigEndian, b.BatchID); err != nil {
		return nil, err
	}

	// Price
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Price)); err != nil {
		return nil, err
	}

	// Amount
	if err := binary.Write(data, binary.BigEndian, PaddingToUint256(b.Amount)); err != nil {
		return nil, err
	}

	// Side
	if err := binary.Write(data, binary.BigEndian, b.Side); err != nil {
		return nil, err
	}

	// Number of orders
	if err := binary.Write(data, binary.BigEndian, uint8(len(b.Orders))); err != nil {
		return nil, err
	}

	// Expanded Orders
	if err := binary.Write(data, binary.BigEndian, ordersData.Bytes()); err != nil {
		return nil, err
	}

	// Owner
	if err := binary.Write(data, binary.BigEndian, b.Owner); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (b *Batch) Sign(_prvkey *ecdsa.PrivateKey) error {
	data, err := b.Encode_Sign()
	if err != nil {
		return err
	}

	// Hash the encode packed data
	hasheddata := crypto.Keccak256Hash(data)

	// Sign the batch
	sig, err := crypto.Sign(hasheddata.Bytes(), _prvkey)
	if err != nil {
		log.Fatal(err)
	}
	b.Signature = sig

	return nil
}

// TODO: batching orders having the same price
func (m *Matcher) batching() []*Batch {
	_logger.Debug("batching...\n")
	batches := []*Batch{}

	m.Mux.Lock()
	defer m.Mux.Unlock()

	threshold := len(m.BidOrders) / 4
	for len(m.BidOrders) > threshold {
		ord := m.BidOrders[0]
		m.BidOrders = m.BidOrders[1:]
		orders := []*MatcherOrder{ord}
		price := ord.Data.Price
		amount := ord.Data.Amount
		for len(m.BidOrders) > 0 && m.BidOrders[0].Data.Price.Cmp(price) == 0 {
			_ord := m.BidOrders[0]
			m.BidOrders = m.BidOrders[1:]
			orders = append(orders, _ord)
			amount = new(big.Int).Add(amount, _ord.Data.Amount)
		}
		batch := m.NewBatch(ord.Data.Price, amount, ord.Data.Side, orders)
		batch.Sign(m.PrivateKey)
		batches = append(batches, batch)
	}

	threshold = len(m.AskOrders) / 4
	for len(m.AskOrders) > threshold {
		ord := m.AskOrders[0]
		m.AskOrders = m.AskOrders[1:]
		orders := []*MatcherOrder{ord}
		price := ord.Data.Price
		amount := ord.Data.Amount
		for len(m.AskOrders) > 0 && m.AskOrders[0].Data.Price.Cmp(price) == 0 {
			_ord := m.AskOrders[0]
			m.AskOrders = m.AskOrders[1:]
			orders = append(orders, _ord)
			amount = new(big.Int).Add(amount, _ord.Data.Amount)
		}
		batch := m.NewBatch(ord.Data.Price, amount, ord.Data.Side, orders)
		batch.Sign(m.PrivateKey)
		batches = append(batches, batch)
	}

	_logger.Debug("got %v batch(es)\n", len(batches))
	return batches
}

func PaddingToUint256(num *big.Int) []byte {
	return append(make([]byte, 32-len(num.Bytes())), num.Bytes()...)
}

func (b *Batch) Equal(_b *Batch) bool {
	if len(b.Orders) != len(_b.Orders) {
		return false
	}
	for idx, order := range b.Orders {
		if !order.Equal(_b.Orders[idx]) {
			return false
		}
	}

	return b.BatchID == _b.BatchID &&
		b.Price.Cmp(_b.Price) == 0 &&
		b.Amount.Cmp(_b.Amount) == 0 &&
		b.Side == _b.Side &&
		b.Owner.Cmp(_b.Owner) == 0 &&
		bytes.Equal(b.Signature, _b.Signature)
}
