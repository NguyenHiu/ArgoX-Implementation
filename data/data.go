package data

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

const DATA_PATH = "data/data.json"

func loadData() (map[string]string, error) {
	file, err := os.Open(DATA_PATH)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]string)
	dec := json.NewDecoder(file)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func saveData(data map[string]string) error {
	file, err := os.Create(DATA_PATH)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", DATA_PATH)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(&data); err != nil {
		return fmt.Errorf("cannot write data into file: %v", DATA_PATH)
	}

	return nil
}

func Get(contractName string) (common.Address, error) {
	data, err := loadData()
	if err != nil {
		return common.Address{}, err
	}

	value, ok := data[contractName]
	if !ok {
		return common.Address{}, fmt.Errorf("cannot find contract name: %v", contractName)
	}

	return common.HexToAddress(value), nil
}

func Set(contractName string, address string) {
	if _, err := os.Stat(DATA_PATH); os.IsNotExist(err) {
		_, err = os.Create(DATA_PATH)
		if err != nil {
			fmt.Printf("Cannot create file %v in Set() func\n", DATA_PATH)
			return
		}
	}

	data, err := loadData()
	if err != nil {
		fmt.Printf("Cannot load data in file %v in Set() func\n", DATA_PATH)
		return
	}

	data[contractName] = address
	if err := saveData(data); err != nil {
		fmt.Printf("Cannot save data in file %v\n", DATA_PATH)
		return
	}
}

func SetMap(data map[string]string) {
	if err := saveData(data); err != nil {
		fmt.Printf("Cannot save data in file %v\n", DATA_PATH)
	}
}
