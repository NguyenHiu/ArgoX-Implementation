package constants

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/verifierApp"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"
)

func DeployContracts(nodeURL string, chainID uint64, privatekey string) (adj common.Address, ahs []common.Address, app common.Address) {
	k, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)
	cb, err := client.CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		panic(err)
	}

	acc := accounts.Account{Address: crypto.PubkeyToAddress(k.PublicKey)}

	// Deploy adjudicator
	adj, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc)
	if err != nil {
		panic(err)
	}

	ahs = []common.Address{}
	// Deploy asset holder
	ah, err := ethchannel.DeployETHAssetholder(context.TODO(), cb, adj, acc)
	if err != nil {
		panic(err)
	}
	ahs = append(ahs, ah)
	// Deploy Gavin asset holder
	ga, err := ethchannel.DeployERC20Assetholder(context.TODO(), cb, adj, common.HexToAddress(constants.GAVIN_TOKEN_ADDRESS), acc)
	if err != nil {
		panic(err)
	}
	ahs = append(ahs, ga)

	// Create a transactor
	const gasLimit = 1100000
	tops, err := cb.NewTransactor(context.TODO(), gasLimit, acc)
	if err != nil {
		panic(err)
	}

	// Deploy Verifier App
	app, tx, _, err := verifierApp.DeployVerifierApp(tops, cb)
	if err != nil {
		panic(err)
	}

	// Waiting for deployment
	_, err = bind.WaitDeployed(context.TODO(), cb, tx)
	if err != nil {
		panic(err)
	}

	return adj, ahs, app
}

func SetupClient(
	bus wire.Bus,
	nodeURL string,
	adjudicator common.Address,
	assets []ethwallet.Address,
	privateKey string,
	app *app.VerifyApp,
	stakes []channel.Bal,
	useTrigger bool,
) *client.AppClient {
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}

	w := swallet.NewWallet(k)
	acc := crypto.PubkeyToAddress(k.PublicKey)

	c, err := client.SetupAppClient(
		bus,
		w,
		acc,
		nodeURL,
		constants.CHAIN_ID,
		adjudicator,
		assets,
		app,
		stakes,
		useTrigger,
	)
	if err != nil {
		panic(err)
	}

	return c
}

// balanceLogger is a utility for logging client balances.
type balanceLogger struct {
	ethClient *ethclient.Client
}

// newBalanceLogger creates a new balance logger for the specified ledger.
func NewBalanceLogger(chainURL string) balanceLogger {
	c, err := ethclient.Dial(chainURL)
	if err != nil {
		panic(err)
	}
	return balanceLogger{ethClient: c}
}

// LogBalances prints the balances of the specified clients.
func (l balanceLogger) LogBalances(clients ...*client.AppClient) {
	bals := make([]*big.Float, len(clients))
	for i, c := range clients {
		bal, err := l.ethClient.BalanceAt(context.TODO(), c.WalletAddress(), nil)
		if err != nil {
			log.Fatal(err)
		}
		bals[i] = client.WeiToEth(bal)
	}
	log.Println("Client balances (ETH):", bals)
}
