package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type OrderData struct {
	Price  int
	Amount int
	Side   bool
}

// random n orders
func RandomOrders(minPrice int, maxPrice int, minAmount int, maxAmount int, n int, saveTo string) []*OrderData {
	orders := []*OrderData{}
	for i := 0; i < n; i++ {
		orders = append(orders, &OrderData{
			Price:  rand.Intn(maxPrice-minPrice+1) + minPrice,
			Amount: rand.Intn(maxAmount-minAmount+1) + minAmount,
			Side:   rand.Intn(2) == 1,
		})
	}

	if saveTo != "" {
		var file *os.File
		file, err := os.Create(saveTo)
		if err != nil {
			fmt.Printf("Cannot create file %v\n", saveTo)
		}

		enc := json.NewEncoder(file)
		if err := enc.Encode(&orders); err != nil {
			fmt.Printf("cannot write data into file: %v\n", saveTo)
		}
	}

	return orders
}
