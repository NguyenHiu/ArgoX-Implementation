package app

import (
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

// Exist order 00000000-0000-0000-0000-000000000000
func (d *VerifyAppData) CheckFinal() bool {
	_, ok := d.OrdersMapping[uuid.UUID{}]
	return ok
}

func (d *VerifyAppData) computeFinalBalances(bals channel.Balances) channel.Balances {
	matcherReceivedETH := &big.Int{}
	matcherReceivedGAV := &big.Int{}

	for _, v := range d.Trades {
		priceETH := v.Price
		amountGVN := v.Amount
		_, ok := d.OrdersMapping[v.BidOrder]
		if ok {
			matcherReceivedETH = new(big.Int).Sub(matcherReceivedETH, new(big.Int).Mul(priceETH, amountGVN))
			matcherReceivedGAV = new(big.Int).Add(matcherReceivedGAV, amountGVN)
		}
		_, ok = d.OrdersMapping[v.AskOrder]
		if ok {
			matcherReceivedETH = new(big.Int).Add(matcherReceivedETH, new(big.Int).Mul(priceETH, amountGVN))
			matcherReceivedGAV = new(big.Int).Sub(matcherReceivedGAV, amountGVN)
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
