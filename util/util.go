package util

import (
	"context"
	"log"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/token"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/verifierApp"
	"github.com/NguyenHiu/lightning-exchange/orderClient"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
)

func DeployPerunContracts(nodeURL string, chainID uint64, privatekey string, gavTokenAddr common.Address) (adj common.Address, ahs []common.Address, app common.Address) {
	k, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)
	cb, err := orderClient.CreateContractBackend(nodeURL, chainID, w)
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
	ga, err := ethchannel.DeployERC20Assetholder(context.TODO(), cb, adj, gavTokenAddr, acc)
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

func DeployCustomSC(nodeURL string, chainID uint64, prvkey string) (common.Address, common.Address) {
	privateKey, err := crypto.HexToECDSA(prvkey)
	if err != nil {
		log.Fatal(err)
	}

	// keystore := keystore.NewKeyStore()
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainID)))
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	token, _, tokenInstance, err := token.DeployToken(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	onchain, _, _, err := onchain.DeployOnchain(auth, client, token)
	if err != nil {
		log.Fatal(err)
	}

	// Mint gavin token
	mintGavinToken(tokenInstance, onchain, client, constants.KEY_DEPLOYER)
	mintGavinToken(tokenInstance, onchain, client, constants.KEY_MATCHER_1)
	mintGavinToken(tokenInstance, onchain, client, constants.KEY_MATCHER_2)
	mintGavinToken(tokenInstance, onchain, client, constants.KEY_ALICE)
	mintGavinToken(tokenInstance, onchain, client, constants.KEY_BOB)

	return token, onchain
}

// // balanceLogger is a utility for logging client balances.
// type balanceLogger struct {
// 	ethClient *ethclient.Client
// }

// // newBalanceLogger creates a new balance logger for the specified ledger.
// func NewBalanceLogger(chainURL string) balanceLogger {
// 	c, err := ethclient.Dial(chainURL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return balanceLogger{ethClient: c}
// }

// // LogBalances prints the balances of the specified clients.
// func (l balanceLogger) LogBalances(clients ...*client.AppClient) {
// 	bals := make([]*big.Float, len(clients))
// 	for i, c := range clients {
// 		bal, err := l.ethClient.BalanceAt(context.TODO(), c.WalletAddress(), nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		bals[i] = client.WeiToEth(bal)
// 	}
// 	log.Println("Client balances (ETH):", bals)
// }

// The contract uses Openzeppelin smart contract to verify ECDSA.
// The Openzeppelin ECDSA contract uses V = {27, 28}
// However, `crypto` packages in Golang uses V = {0, 1}
// Therefore, we have to manually change the V value
func CorrectSignToOnchain(signature []byte) []byte {
	if signature[64] == 0 || signature[64] == 1 {
		_signature := make([]byte, len(signature))
		copy(_signature, signature)
		_signature[64] += 27
		return _signature
	}

	return signature

}

func CorrectSignToBE(signature []byte) []byte {
	if signature[64] == 27 || signature[64] == 28 {
		_signature := make([]byte, len(signature))
		copy(_signature, signature)
		_signature[64] -= 27
		return _signature
	}

	return signature
}

func PaddingToUint256(num int64) []byte {
	bigInt := big.NewInt(num)
	return append(make([]byte, 32-len(bigInt.Bytes())), bigInt.Bytes()...)
}

func mintGavinToken(tokenInstance *token.Token, onchainAddr common.Address, client *ethclient.Client, privateKeyHex string) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}

	addr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Mint token
	prepareNonceAndGasPrice(auth, client, addr)
	if _, err = tokenInstance.Mint(auth, addr, big.NewInt(constants.NO_MINTED_GVN_TOKEN)); err != nil {
		log.Fatal(err)
	}

	// Approve onchain contract
	prepareNonceAndGasPrice(auth, client, addr)
	if _, err := tokenInstance.Approve(auth, onchainAddr, big.NewInt(constants.NO_GVN_APPROVE)); err != nil {
		log.Fatal(err)
	}
}

func DepositETH(onchainInstance *onchain.Onchain, _client *ethclient.Client, privateKeyHex string) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}

	addr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Deposit ETH to the exchange
	prepareNonceAndGasPrice(auth, _client, addr)
	auth.Value = big.NewInt(constants.NO_ETH_DEPOSIT)
	// auth.Value = client.EthToWei(big.NewFloat(100))
	if _, err := onchainInstance.MyDeposit(auth); err != nil {
		log.Fatal(err)
	}
}

func prepareNonceAndGasPrice(auth *bind.TransactOpts, client *ethclient.Client, address common.Address) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.Value = &big.Int{}
	auth.GasLimit = uint64(300000)
}

func CalculateTotalUsedGas(addr common.Address) int {
	totalGas := 0
	_client, _ := ethclient.Dial(constants.CHAIN_URL)

	for i := new(big.Int); ; i.Add(i, big.NewInt(1)) {
		block, err := _client.BlockByNumber(context.Background(), i)
		if err != nil {
			return totalGas
		}
		for _, tx := range block.Transactions() {
			if from, err := types.Sender(types.NewLondonSigner(big.NewInt(1337)), tx); err == nil {
				if from.Cmp(addr) == 0 {
					receipt, err := _client.TransactionReceipt(context.Background(), tx.Hash())
					if err != nil {
						log.Fatal(err)
					}
					totalGas += int(receipt.GasUsed)
				}
			}
		}

	}
}
