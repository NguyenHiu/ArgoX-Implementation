package deploy

import (
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/util"
)

func DeployContracts() {
	contractsMap := make(map[string]string)
	token, onchain := util.DeployCustomSC(constants.CHAIN_URL, uint64(constants.CHAIN_ID), constants.KEY_DEPLOYER)
	adj, assetHolders, appAddr := util.DeployPerunContracts(constants.CHAIN_URL, uint64(constants.CHAIN_ID), constants.KEY_DEPLOYER, token)

	contractsMap["token"] = token.String()
	contractsMap["onchain"] = onchain.String()
	contractsMap["adj"] = adj.String()
	contractsMap["ethholder"] = assetHolders[constants.ETH].String()
	contractsMap["gvnholder"] = assetHolders[constants.GVN].String()
	contractsMap["appaddr"] = appAddr.String()

	data.SetMap(contractsMap)
}
