package worker

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"sync"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

var _logger = logger.NewLogger("worker", logger.None, logger.None)

type Worker struct {
	BidBatches []*Batch
	AskBatches []*Batch
	Batches    map[uuid.UUID]*Batch

	OnchainIstance *onchain.Onchain
	PrivateKey     *ecdsa.PrivateKey
	WatchOpts      *bind.WatchOpts
	Client         *ethclient.Client
	Address        common.Address
	Auth           *bind.TransactOpts

	Mux sync.Mutex
}

func NewWorker(onchainAddr common.Address, privateKey string, client *ethclient.Client) *Worker {
	_privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	onchainInstance, err := onchain.NewOnchain(onchainAddr, client)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(_privateKey, big.NewInt(constants.CHAIN_ID))
	if err != nil {
		log.Fatal(err)
	}

	newWorker := &Worker{
		BidBatches:     make([]*Batch, 0),
		AskBatches:     make([]*Batch, 0),
		Batches:        make(map[uuid.UUID]*Batch),
		PrivateKey:     _privateKey,
		OnchainIstance: onchainInstance,
		WatchOpts:      &bind.WatchOpts{Context: context.Background()},
		Client:         client,
		Address:        crypto.PubkeyToAddress(_privateKey.PublicKey),
		Auth:           auth,
	}

	return newWorker
}
