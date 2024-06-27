package app

import (
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

// Exist order 00000000-0000-0000-0000-000000000000
func (d *VerifyAppData) CheckFinal() bool {
	_, ok := d.Orders[uuid.UUID{}]
	return ok
}

func (d *VerifyAppData) computeFinalBalances(bals channel.Balances) channel.Balances {
	matcherReceivedETH := &big.Int{}
	matcherReceivedGAV := &big.Int{}

	for k, v := range d.Orders {
		if len(d.Msgs[k]) != 0 {
			if v.Side == constants.BID {
				matcherReceivedETH = new(big.Int).Add(matcherReceivedETH, v.Price)
				matcherReceivedGAV = new(big.Int).Sub(matcherReceivedGAV, d.Msgs[k][len(d.Msgs[k])-1].MatchedAmount)
			} else {
				matcherReceivedETH = new(big.Int).Sub(matcherReceivedETH, v.Price)
				matcherReceivedGAV = new(big.Int).Add(matcherReceivedGAV, d.Msgs[k][len(d.Msgs[k])-1].MatchedAmount)
			}
		}
	}

	finalBals := bals.Clone()
	finalBals[constants.ETH][constants.MATCHER] = new(big.Int).Add(bals[constants.ETH][constants.MATCHER], matcherReceivedETH)
	finalBals[constants.ETH][constants.TRADER] = new(big.Int).Sub(bals[constants.ETH][constants.TRADER], matcherReceivedETH)
	finalBals[constants.GVN][constants.MATCHER] = new(big.Int).Add(bals[constants.GVN][constants.MATCHER], matcherReceivedGAV)
	finalBals[constants.GVN][constants.TRADER] = new(big.Int).Sub(bals[constants.GVN][constants.TRADER], matcherReceivedGAV)

	return finalBals
}

func PaddingToUint256(num *big.Int) []byte {
	return append(make([]byte, 32-len(num.Bytes())), num.Bytes()...)
}
