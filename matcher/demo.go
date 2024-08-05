package matcher

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"

	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

func (m *Matcher) Matching() {
	m.matching()
}

func (m *Matcher) Batching(noBatch int) {
	if len(m.BidOrders) >= noBatch || len(m.AskOrders) >= noBatch {
		batches := m.batching()
		m.BatchesList = append(m.BatchesList, batches...)
		for _, batch := range batches {
			m.BatchStatusMapping[batch.BatchID] = CREATED
		}
	}
}

type OrderStatus struct {
	ID     uuid.UUID
	Status string
}

func (m *Matcher) GetOrderStatus(orderID uuid.UUID) string {
	if status, ok := m.OrderStatusMapping[orderID]; ok {
		return status
	}
	return ""
}

type BatchStatus struct {
	ID     uuid.UUID
	Status string
}

func (m *Matcher) WatchingDemo(instance *onchain.Onchain) {
	opt := bind.WatchOpts{Context: context.Background()}
	go m.WatchAcceptBatch(instance, &opt)
	go m.WatchReceivedBatchDetails(instance, &opt)
	go m.WatchRemoveBatchOutOfDate(instance, &opt)
	go m.WatchValidMatching(instance, &opt)
}

func (m *Matcher) WatchAcceptBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainAcceptBatch)
	sub, err := instance.WatchAcceptBatch(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			if _, ok := m.Batches[id]; ok {
				m.BatchStatusMapping[id] = ACCEPTED
			}
		}
	}
}

func (m *Matcher) WatchReceivedBatchDetails(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainReceivedBatchDetails)
	sub, err := instance.WatchReceivedBatchDetails(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			if _, ok := m.Batches[id]; ok {
				m.Mux.Lock()
				if _status, ok := m.BatchStatusMapping[id]; ok {
					if _status != FULFILED {
						m.BatchStatusMapping[id] = WAITING_FOR_OTHER
					}
					m.Mux.Unlock()
				}
			}

		}
	}
}

func (m *Matcher) WatchRemoveBatchOutOfDate(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainRemoveBatchOutOfDate)
	sub, err := instance.WatchRemoveBatchOutOfDate(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			if _, ok := m.Batches[id]; ok {
				m.BatchStatusMapping[id] = REPORTED
			}
		}
	}
}

func (m *Matcher) SendBatchDetails(batchID uuid.UUID) {
	m.Mux.Lock()
	batch, ok := m.Batches[batchID]
	m.Mux.Unlock()
	if !ok {
		return
	}
	onchainOrders := []onchain.OnchainOrder{}
	for _, order := range batch.Orders {
		trades := new(bytes.Buffer)
		for _, trade := range order.Trades {
			tData, err := trade.Encode_TransferBatching()
			if err != nil {
				_logger.Error("sending batch's detail got error, encode trade, err: %v\n", err)
			}
			if err := binary.Write(trades, binary.BigEndian, tData); err != nil {
				_logger.Error("sending batch's detail, hash a trade got error, err: %v\n", err)
			}
		}
		tradeHash := crypto.Keccak256Hash(trades.Bytes())

		oData, err := order.OriginalOrder.Encode_TransferBatching()
		if err != nil {
			_logger.Error("sending batch's detail, hash original order got error, err: %v\n", err)
		}
		originalOrderHash := crypto.Keccak256Hash(oData)

		onchainOrders = append(onchainOrders, onchain.OnchainOrder{
			Price:             order.ShadowOrder.Price,
			Amount:            order.ShadowOrder.Amount,
			Side:              order.ShadowOrder.Side,
			From:              order.ShadowOrder.From,
			TradeHash:         tradeHash,
			OriginalOrderHash: originalOrderHash,
			Owner:             common.Address(*order.OriginalOrder.Owner),
		})
	}

	// Send batch's details
	// m.Mux.Lock()
	m.prepareNonceAndGasPrice(0, 900000)
	_, err := m.OnchainInstance.SubmitOrderDetails(m.Auth, batchID, onchainOrders)
	if err != nil {
		_logger.Error("Submit Error, err: %v\n", err)
	}
	_logger.Debug("Matcher::%v, Submit batch's details, batch::%v\n", m.ID.String()[:6], batchID.String())
	// m.Mux.Unlock()
}

func (m *Matcher) WatchValidMatching(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainValidMatching)
	sub, err := instance.WatchValidMatching(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			if _, ok := m.Batches[id]; ok {
				_logger.Debug("ValidMatching, got id: %v\n", id.String())
				m.Mux.Lock()
				m.BatchStatusMapping[id] = FULFILED
				_logger.Debug("batch::%v - Status: %v\n", id, m.BatchStatusMapping[id])
				m.Mux.Unlock()
			}
		}
	}
}

func (m *Matcher) WatchRevertBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainRevertBatch)
	sub, err := instance.WatchRevertBatch(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			id, _ := uuid.FromBytes(vLogs.Arg0[:])
			if _, ok := m.Batches[id]; ok {
				_logger.Debug("Revert Batch, got id: %v\n", id.String())
				m.Mux.Lock()
				m.BatchStatusMapping[id] = REVERTED
				_logger.Debug("batch::%v - Status: %v\n", id, m.BatchStatusMapping[id])
				m.Mux.Unlock()
			}
		}
	}
}
