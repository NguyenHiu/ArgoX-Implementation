package data

import "math/rand"

type OrderData struct {
	Price  int
	Amount int
	Side   bool
}

// random n orders
func RandomOrders(minPrice int, maxPrice int, minAmount int, maxAmount int, n int) []OrderData {
	orders := []OrderData{}
	for i := 0; i < n; i++ {
		orders = append(orders, OrderData{
			Price:  rand.Intn(maxPrice-minPrice+1) + minPrice,
			Amount: rand.Intn(maxAmount-minAmount+1) + minAmount,
			Side:   rand.Intn(2) == 1,
		})
	}

	return orders
}
