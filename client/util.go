package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
)

const (
	txFinalityDepth = 1 // Number of blocks required to confirm a transaction.
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

	return ethchannel.NewContractBackend(ethClient, transactor, txFinalityDepth), nil
}
