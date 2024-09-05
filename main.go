package main

import (
	"log"
	"time"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/deploy"
	"github.com/NguyenHiu/lightning-exchange/listener"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/reporter"
	"github.com/NguyenHiu/lightning-exchange/server"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/NguyenHiu/lightning-exchange/worker"
	"github.com/ethereum/go-ethereum/ethclient"
)

var _logger = logger.NewLogger("Main", logger.Red, logger.Bold)

func main() {
	constants.CHAIN_URL = "ws://127.0.0.1:8545"
	constants.CHAIN_ID = 1337
	constants.NO_MATCHER = 2
	constants.SEND_TO = 1

	// Get ethclient
	clientNode, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy contracts
	_logger.Debug("Deploy smart contracts...\n")
	deploy.DeployContracts()
	_token, _onchain, adj, assetHolders, appAddr := getContracts()

	// Deposit ETH to the smart contract
	_logger.Debug("Deposit ETH to the exchange...\n")
	onchainInstance, _ := onchain.NewOnchain(_onchain, clientNode)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_ALICE)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_BOB)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_DEPLOYER)

	// Listener
	_logger.Debug("Start Listener...\n")
	_listenerInstance := listener.NewListener()
	go _listenerInstance.StartListener(_onchain)

	// Reporter
	_logger.Debug("Start Reporter...\n")
	rp, err := reporter.NewReporter(_onchain, constants.KEY_REPORTER, constants.CHAIN_ID)
	if err != nil {
		_logger.Error("Create reporter error, err: %v\n", err)
	}
	rp.Listening()

	// Worker
	_logger.Debug("Start Worker...\n")
	w := worker.NewWorker(_onchain, constants.KEY_WORKER, clientNode)
	w.Listening()

	// Super Matcher
	_logger.Debug("Start Super Matcher...\n")
	sm := SetupSuperMatcher(_onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	// Slice of matchers
	_logger.Debug("Init matchers...\n")
	matchers := []*matcher.Matcher{}
	for i := 0; i < constants.NO_MATCHER; i++ {
		matcher := matcher.NewMatcher(
			assetHolders,
			adj,
			appAddr,
			_onchain,
			constants.MATCHER_PRVKEYS[i],
			clientNode,
			constants.CHAIN_ID,
			_token,
			sm,
		)
		matcher.Register()
		go matcher.WatchingDemo(onchainInstance)
		matchers = append(matchers, matcher)
	}

	// Alice: Initialize connections to all matchers
	alice := user.NewUser(constants.KEY_ALICE)
	for i := 0; i < constants.NO_MATCHER; i++ {
		orderBus, tradeBus := matchers[i].SetupClient(alice.ID)
		alice.SetupClient(
			matchers[i].ID, orderBus,
			tradeBus,
			constants.CHAIN_URL,
			matchers[i].Adjudicator,
			matchers[i].AssetHolders,
			matchers[i].OrderApp,
			matchers[i].TradeApp,
			matchers[i].Stakes,
			matchers[i].EmptyStakes,
			_token,
		)

		if ok := matchers[i].OpenAppChannel(
			alice.ID,
			alice.Connections[matchers[i].ID].OrderAppClient.WireAddress(),
		); !ok {
			log.Fatalln("OpenAppChannel Failed")
		}
		alice.AcceptedChannel(matchers[i].ID)
		<-time.After(time.Millisecond * 500)
	}

	// Onchain
	_listenerInstance.IsGetPriceCurve = true
	go _listenerInstance.GetPriceCurve()

	_server := server.NewServer(
		7000,
		alice,
		matchers,
		sm,
		rp,
		_listenerInstance,
		w,
	)
	_server.Start()
}
