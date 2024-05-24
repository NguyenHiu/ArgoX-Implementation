#!/bin/sh

### Make sure you have already installed `abigen` and `solc`
# Generate `bin` and `abi` files
solc --bin VerifierApp.sol -o build
solc --abi VerifierApp.sol -o build

# Create `generated` folder
mkdir -p generated/verifierApp

# Generate smart contract as `go` file
abigen --bin=./build/VerifierApp.bin --abi=./build/VerifierApp.abi --pkg=verifierApp --out=generated/verifierApp/verifierApp.go