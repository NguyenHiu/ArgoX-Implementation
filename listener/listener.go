package listener

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

var _logger = logger.NewLogger("Listener", logger.Yellow, logger.Bold)

// TODO: Collect Time & Matched Amount & Price Onchain
type Listener struct {
	TotalTimeOnchain          int
	TotalMatchedAmountOnchain *big.Int
	NumberOfMatchedOrder      int64
	PriceCurveOnchain         []*big.Int
	CurrentPrice              *big.Int
	IsGetPriceCurve           bool
	TotalProfitOnchain        *big.Int
	TotalRawProfitOnchain     *big.Int

	batchPriceMapping map[uuid.UUID]*big.Int
}

func NewListener() *Listener {
	return &Listener{
		TotalTimeOnchain:          0,
		TotalMatchedAmountOnchain: new(big.Int),
		NumberOfMatchedOrder:      0,
		PriceCurveOnchain:         []*big.Int{},
		CurrentPrice:              new(big.Int),
		IsGetPriceCurve:           false,
		TotalProfitOnchain:        new(big.Int),
		TotalRawProfitOnchain:     new(big.Int),
		batchPriceMapping:         make(map[uuid.UUID]*big.Int),
	}
}

func (l *Listener) StartListener(onchainAddr common.Address) {
	client, _ := ethclient.Dial(constants.CHAIN_URL)
	instance, _ := onchain.NewOnchain(onchainAddr, client)

	opts := bind.WatchOpts{Context: context.Background()}
	go l.WatchFullfilMatch(instance, &opts)
	go l.WatchReceivedBatchDetails(instance, &opts)
	go l.WatchAcceptBatch(instance, &opts)
	go l.WatchPunishMatcher(instance, &opts)
	go l.WatchRemoveBatchOutOfDate(instance, &opts)
	go l.WatchInvalidBatch(instance, &opts)
	go l.WatchInvalidOrder(instance, &opts)
	go l.WatchRevertBatch(instance, &opts)
	go l.WatchLogString(instance, &opts)
	go l.WatchLogBytes32(instance, &opts)
	go l.WatchLogBytes16(instance, &opts)
	go l.WatchLogAddress(instance, &opts)
	go l.WatchLogBytes(instance, &opts)
	go l.WatchLogUint256(instance, &opts)
	go l.WatchLogRecoverError(instance, &opts)

	// Statistical
	go l.WatchLogMatchingTimestamp(instance, &opts)
	go l.WatchMatchPrice(instance, &opts)
	go l.WatchMatchAmount(instance, &opts)
	go l.WatchBatchRawProfit(instance, &opts)
}

// Statistical
func (l *Listener) GetPriceCurve() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if l.IsGetPriceCurve {
			l.PriceCurveOnchain = append(
				l.PriceCurveOnchain,
				l.CurrentPrice,
			)
		}
	}
}

// Statistical
func (l *Listener) WatchMatchAmount(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainBatchMatchAmountAndProfit)
	sub, err := instance.WatchBatchMatchAmountAndProfit(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			l.TotalMatchedAmountOnchain.Add(
				l.TotalMatchedAmountOnchain,
				vLogs.Arg0,
			)
			l.TotalMatchedAmountOnchain.Add(
				l.TotalMatchedAmountOnchain,
				vLogs.Arg0,
			)
			vLogs.Arg1.Mul(vLogs.Arg1, vLogs.Arg0)
			vLogs.Arg1.Mul(vLogs.Arg1, big.NewInt(2))
			l.TotalProfitOnchain.Add(
				l.TotalProfitOnchain,
				vLogs.Arg1,
			)
		}
	}
}

// Statistical
func (l *Listener) WatchMatchPrice(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainMatchedPrice)
	sub, err := instance.WatchMatchedPrice(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			l.CurrentPrice = vLogs.Arg0
		}
	}
}

// Statistical
func (l *Listener) WatchBatchRawProfit(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainBatchRawProfit)
	sub, err := instance.WatchBatchRawProfit(opt, logs)
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
			l.batchPriceMapping[id] = vLogs.Arg1
		}
	}
}

// Statistical
func (l *Listener) WatchLogMatchingTimestamp(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainBatchTimestamp)
	sub, err := instance.WatchBatchTimestamp(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	batchTimestampMapping := make(map[uuid.UUID]*big.Int, 0)
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_batchID, _ := uuid.FromBytes(vLogs.Arg0[:])
			if _startTime, _ok := batchTimestampMapping[_batchID]; _ok {
				l.TotalTimeOnchain += int(vLogs.Arg1.Sub(vLogs.Arg1, _startTime).Int64())
			} else {
				batchTimestampMapping[_batchID] = vLogs.Arg1
			}
		}
	}
}

func (l *Listener) WatchFullfilMatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainFullfilMatch)
	sub, err := instance.WatchFullfilMatch(opt, logs)
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
			_logger.Info("[Fullfill] Batch::%v\n", id.String())
		}
	}
}

func (l *Listener) WatchReceivedBatchDetails(instance *onchain.Onchain, opt *bind.WatchOpts) {
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
			_logger.Info("[Details] Batch::%v\n", id.String())
			l.NumberOfMatchedOrder += vLogs.Arg1.Int64()

			if price, ok := l.batchPriceMapping[id]; !ok {
				log.Fatal("Price not found")
			} else {
				l.TotalRawProfitOnchain.Add(
					l.TotalRawProfitOnchain,
					new(big.Int).Mul(price, vLogs.Arg1),
				)
			}

		}
	}
}

func (l *Listener) WatchAcceptBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
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
			_logger.Info("[Accept] Batch::%v\n", id.String())
			// LogOrderBookOverview(instance)
		}
	}
}

func (l *Listener) WatchPunishMatcher(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainPunishMatcher)
	sub, err := instance.WatchPunishMatcher(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Info("[Punish] Matcher::%v\n", vLogs.Arg0.String())
		}
	}
}

func (l *Listener) WatchRemoveBatchOutOfDate(instance *onchain.Onchain, opt *bind.WatchOpts) {
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
			_logger.Info("[Remove] Batch::%v\n", id.String())
		}
	}
}

func (l *Listener) WatchInvalidOrder(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainInvalidOrder)
	sub, err := instance.WatchInvalidOrder(opt, logs)
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
			_logger.Info("[Invalid Order] Batch::%v\n", id.String())
		}
	}
}

func (l *Listener) WatchInvalidBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainInvalidBatch)
	sub, err := instance.WatchInvalidBatch(opt, logs)
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
			_logger.Info("[Invalid Batch] Batch::%v\n", id.String())
		}
	}
}

func (l *Listener) WatchRevertBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
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
			_logger.Info("[Revert] Batch::%v\n", id.String())
			log.Fatal("Revert batch")
			// LogOrderBookOverview(instance)
		}
	}
}

func (l *Listener) WatchLogString(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogString)
	sub, err := instance.WatchLogString(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Debug("[Contract] %v\n", vLogs.Arg0)
		}
	}
}

func (l *Listener) WatchLogAddress(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogAddress)
	sub, err := instance.WatchLogAddress(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Debug("[Contract] %v\n", vLogs.Arg0)
		}
	}
}

func (l *Listener) WatchLogBytes32(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogBytes32)
	sub, err := instance.WatchLogBytes32(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Debug("[Contract] %v\n", vLogs.Arg0)
		}
	}
}

func (l *Listener) WatchLogBytes16(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogBytes16)
	sub, err := instance.WatchLogBytes16(opt, logs)
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
			_logger.Debug("[Contract] %v\n", id.String())
		}
	}
}

func (l *Listener) WatchLogBytes(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogBytes)
	sub, err := instance.WatchLogBytes(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Debug("[Contract] %v\n", vLogs.Arg0)
		}
	}
}

func (l *Listener) WatchLogUint256(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogUint256)
	sub, err := instance.WatchLogUint256(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Debug("[Contract] %v\n", vLogs.Arg0)
		}
	}
}

func (l *Listener) WatchLogRecoverError(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainLogRecoverError)
	sub, err := instance.WatchLogRecoverError(opt, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLogs := <-logs:
			_logger.Debug("[Contract] %v\n", vLogs.Arg0)
		}
	}
}
