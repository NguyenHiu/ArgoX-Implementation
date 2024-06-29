package listener

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

var _logger = logger.NewLogger("Listener", logger.Yellow, logger.Bold)

func StartListener(onchainAddr common.Address) {
	client, _ := ethclient.Dial(constants.CHAIN_URL)
	instance, _ := onchain.NewOnchain(onchainAddr, client)

	opts := bind.WatchOpts{Context: context.Background()}
	go WatchPartialMatch(instance, &opts)
	go WatchFullfilMatch(instance, &opts)
	go WatchReceivedBatchDetails(instance, &opts)
	go WatchAcceptBatch(instance, &opts)
	go WatchPunishMatcher(instance, &opts)
	go WatchRemoveBatchOutOfDate(instance, &opts)
	go WatchInvalidBatch(instance, &opts)
	go WatchInvalidOrder(instance, &opts)
	go WatchRevertBatch(instance, &opts)
	go WatchLogString(instance, &opts)
	go WatchLogBytes32(instance, &opts)
	go WatchLogBytes16(instance, &opts)
	go WatchLogAddress(instance, &opts)
	go WatchLogBytes(instance, &opts)
	go WatchLogRecoverError(instance, &opts)
}

func WatchPartialMatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainPartialMatch)
	sub, err := instance.WatchPartialMatch(opt, logs)
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
			_logger.Info("[Partial] Batch::%v\n", id.String())
		}
	}
}

func WatchFullfilMatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchReceivedBatchDetails(instance *onchain.Onchain, opt *bind.WatchOpts) {
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
		}
	}
}

func WatchAcceptBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchPunishMatcher(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchRemoveBatchOutOfDate(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchInvalidOrder(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchInvalidBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchRevertBatch(instance *onchain.Onchain, opt *bind.WatchOpts) {
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
			// LogOrderBookOverview(instance)
		}
	}
}

func WatchLogString(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchLogAddress(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchLogBytes32(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchLogBytes16(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchLogBytes(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func WatchLogRecoverError(instance *onchain.Onchain, opt *bind.WatchOpts) {
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

func LogOrderBookOverview(instance *onchain.Onchain) {
	bidBatches, _ := instance.GetBidOrders(&bind.CallOpts{Context: context.Background()})
	askBatches, _ := instance.GetAskOrders(&bind.CallOpts{Context: context.Background()})
	zero := &big.Int{}
	s := ""
	price := big.NewInt(0)
	amount := big.NewInt(0)

	s += "(Ask)\n"
	for i := len(askBatches) - 1; i >= 0; i-- {
		if price.Cmp(zero) == 0 {
			price = askBatches[i].Price
			amount = askBatches[i].Amount
		} else if price.Cmp(askBatches[i].Price) == 0 {
			amount = new(big.Int).Add(amount, askBatches[i].Amount)
		} else {
			s += "\t" + price.String() + ";\t" + amount.String() + "\n"
			price = zero
		}
	}
	if price.Cmp(zero) != 0 {
		s += "\t" + price.String() + ";\t" + amount.String() + "\n"
		price = zero
	}

	s += "\t-----------------------\n"

	s += "(Bid)\n"
	for _, batch := range bidBatches {
		if price.Cmp(zero) == 0 {
			price = batch.Price
			amount = batch.Amount
		} else if price.Cmp(batch.Price) == 0 {
			amount = new(big.Int).Add(amount, batch.Amount)
		} else {
			s += "\t" + price.String() + ";\t" + amount.String() + "\n"
			price = zero
		}
	}
	if price.Cmp(zero) != 0 {
		s += "\t" + price.String() + ";\t" + amount.String() + "\n"
	}

	_logger.Debug("Overview:\n%v\n", s)
}
