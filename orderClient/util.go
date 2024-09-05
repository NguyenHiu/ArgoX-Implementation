package orderClient

import (
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
)

func CreateContractBackend(
	nodeURL string,
	chainID uint64,
	w *swallet.Wallet,
) (ethchannel.ContractBackend, error) {
	signer := types.NewEIP155Signer(new(big.Int).SetUint64(chainID))
	transactor := swallet.NewTransactor(w, signer)

	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return ethchannel.ContractBackend{}, err
	}

	return ethchannel.NewContractBackend(ethClient, transactor, constants.TX_FINALITY_DEPTH), nil
}

// WalletAddress returns the wallet address of the client.
func (c *OrderAppClient) WalletAddress() common.Address {
	return common.Address(*c.account.(*ethwallet.Address))
}

// WalletAddressAsEthwallet returns the wallet address of the client in the format of ethwallet.Address
func (c *OrderAppClient) WalletAddressAsEthwallet() *ethwallet.Address {
	wallet, ok := c.account.(*ethwallet.Address)
	if !ok {
		panic("can not convert")
	}
	return wallet
}

// WireAddress returns the wire address of the client.
func (c *OrderAppClient) WireAddress() *ethwallet.Address {
	return ethwallet.AsWalletAddr(common.Address(c.account.Bytes()))
}

// EthToWei converts a given amount in ETH to Wei.
func EthToWei(ethAmount *big.Float) (weiAmount *big.Int) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).Mul(ethAmount, weiPerEthFloat)
	weiAmount, _ = weiAmountFloat.Int(nil)
	return weiAmount
}

// WeiToEth converts a given amount in Wei to ETH.
func WeiToEth(weiAmount *big.Int) (ethAmount *big.Float) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).SetInt(weiAmount)
	return new(big.Float).Quo(weiAmountFloat, weiPerEthFloat)
}
