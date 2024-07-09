package orderApp

import (
	"math/big"

	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

// Exist order 00000000-0000-0000-0000-000000000000
func (d *OrderAppData) CheckFinal() bool {
	_, ok := d.OrdersMapping[uuid.UUID{}]
	return ok
}

func (d *OrderAppData) computeFinalBalances(bals channel.Balances) channel.Balances {
	return bals
}

func PaddingToUint256(num *big.Int) []byte {
	return append(make([]byte, 32-len(num.Bytes())), num.Bytes()...)
}
