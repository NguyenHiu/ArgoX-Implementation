package supermatcher

import (
	"math/big"

	"github.com/google/uuid"
)

type OrderRequirements struct {
	IsValidSignature bool
}

type BatchRequirements struct {
	IsValidMatcher   bool
	IsValidSignature bool
	IsValidOrders    bool
}

type BatchResult struct {
	BatchID         uuid.UUID
	Orders          []uuid.UUID
	BatchStatus     BatchRequirements
	OrdersStatus    []OrderRequirements
	DuplicateOrders []uuid.UUID
}

func (sm *SuperMatcher) GetRemainingAmount(oderID uuid.UUID) *big.Int {
	if _amount, ok := sm.MatchedOrders[oderID]; ok {
		return _amount
	}
	return big.NewInt(0)
}
