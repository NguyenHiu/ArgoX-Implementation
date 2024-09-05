package reporter

import (
	"context"
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

var _logger = logger.NewLogger("Reporter", logger.Blue, logger.Bold)

type Reporter struct {
	ReporterID      uuid.UUID
	OnchainInstance *onchain.Onchain
	Auth            *bind.TransactOpts
	Client          *ethclient.Client
	Address         common.Address
	WaitingTime     int64
	PendingBatches  map[uuid.UUID]int64
	WatchOpts       *bind.WatchOpts
	Mux             sync.Mutex
}

func NewReporter(onchainAdrr common.Address, prvateKeyHex string, chainID int64) (*Reporter, error) {

	client, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		return nil, err
	}

	instance, _ := onchain.NewOnchain(onchainAdrr, client)

	id, _ := uuid.NewRandom()
	_prvkey, err := crypto.HexToECDSA(prvateKeyHex)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(_prvkey, big.NewInt(chainID))
	if err != nil {
		return nil, err
	}

	reporter := &Reporter{
		ReporterID:      id,
		OnchainInstance: instance,
		Auth:            auth,
		Client:          client,
		Address:         crypto.PubkeyToAddress(_prvkey.PublicKey),
		WaitingTime:     getWaitingTime(instance).Int64(),
		PendingBatches:  make(map[uuid.UUID]int64),
		WatchOpts:       &bind.WatchOpts{Context: context.Background()},
	}

	return reporter, nil
}
