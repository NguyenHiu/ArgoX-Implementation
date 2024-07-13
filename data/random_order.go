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
func RandomOrders(minPrice int, maxPrice int, minAmount int, maxAmount int, n int) []*OrderData {
	orders := []*OrderData{}
	for i := 0; i < n; i++ {
		orders = append(orders, &OrderData{
			Price:  rand.Intn(maxPrice-minPrice+1) + minPrice,
			Amount: rand.Intn(maxAmount-minAmount+1) + minAmount,
			Side:   rand.Intn(2) == 1,
		})
	}

	return orders
}

func LoadOrders(path string) ([]*OrderData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	orders := make([]*OrderData, 0)
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func SaveOrders(orders []*OrderData, filename string) error {
	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Cannot create file %v\n", filename)
		}

		defer file.Close()

		enc := json.NewEncoder(file)
		if err := enc.Encode(&orders); err != nil {
			fmt.Printf("cannot write data into file: %v\n", filename)
		}
	}

	return nil
}
