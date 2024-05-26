#!/bin/sh

### Make sure you have already installed `abigen` and `solc`
# Generate `bin` and `abi` files
solc --bin Token.sol -o build
solc --abi Token.sol -o build

# Create `generated` folder
mkdir -p generated/Token

# Generate smart contract as `go` file
abigen --bin=./build/Token.bin --abi=./build/Token.abi --pkg=build --out=generated/Token/Token.go