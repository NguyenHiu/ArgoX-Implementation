#!/bin/bash

./delete_build.sh

files=$(ls contracts)

file_args=""
for file in $files
do 
    file_args="$file_args /contracts/contracts/$file"
done

echo "> Generating abi, bin files in build folder..."

sudo docker run -v $(pwd):/contracts ethereum/solc:0.8.24 -o /contracts/build/ --overwrite --abi --bin $file_args 

echo "> Done"

echo "> Generting smart contracts as go files..."
for file in $files
do 
    filename=$(basename "$file" .sol)
    packagename=${filename,}
    mkdir -p ./generated/$packagename
    echo "> Generating $filename file..."
    abigen --bin=./build/$filename.bin --abi=./build/$filename.abi --pkg=$packagename --out=./generated/$packagename/$filename.go
    echo "> Done"
done