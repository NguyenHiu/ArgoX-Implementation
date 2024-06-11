#!/bin/sh

rm -rf ./build

### Make sure you have already installed `abigen` and `solc`
# Generate `bin` and `abi` files
solc --bin VerifierApp.sol -o build
solc --abi VerifierApp.sol -o build

rm -rf ./generated/verifierApp
# Create `generated` folder
mkdir -p generated/verifierApp

# Generate smart contract as `go` file
abigen --bin=./build/VerifierApp.bin --abi=./build/VerifierApp.abi --pkg=verifierApp --out=generated/verifierApp/verifierApp.go


# # Generate `bin` and `abi` files
solc --bin Onchain.sol -o build
solc --abi Onchain.sol -o build

rm -rf ./generated/onchain
# Create `generated` folder
mkdir -p generated/onchain

# Generate smart contract as `go` file
abigen --bin=./build/Onchain.bin --abi=./build/Onchain.abi --pkg=onchain --out=generated/onchain/onchain.go


# Generate `bin` and `abi` files
solc --bin Token.sol -o build
solc --abi Token.sol -o build

rm -rf ./generated/token
# Create `generated` folder
mkdir -p generated/token

# #Generate smart contract as `go` file
abigen --bin=./build/Token.bin --abi=./build/Token.abi --pkg=token --out=generated/token/token.go
