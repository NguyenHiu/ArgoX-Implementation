package listener

import (
	"context"
	"log"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

var _logger = logger.NewLogger("Listener")

func StartListener(onchainAddr common.Address) {
	client, _ := ethclient.Dial(constants.CHAIN_URL)
	instance, _ := onchain.NewOnchain(onchainAddr, client)

	opts := bind.WatchOpts{Context: context.Background()}
	go watchFullfilEvent(instance, &opts)
	go watchReceiveEvent(instance, &opts)
}

func watchFullfilEvent(contract *onchain.Onchain, opts *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainFullfilMatch)
	sub, err := contract.WatchFullfilMatch(opts, logs)
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
			_logger.Info("[Fullfill] batch::%v\n", id)
		}
	}
}

func watchReceiveEvent(contract *onchain.Onchain, opts *bind.WatchOpts) {
	logs := make(chan *onchain.OnchainAcceptBatch)
	sub, err := contract.WatchAcceptBatch(opts, logs)
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
			_logger.Info("[Accept] batch::%v\n", id)
		}
	}
}
