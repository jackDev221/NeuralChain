#!/usr/bin/env bash
set -euo pipefail

pushd ./consensus/staking_contracts
# Generate bytes code and abi files
echo "Generate bytes code and abi files"
# Download smart-contact library, it might be long
pushd staking-contract
yarn install
popd
docker run -v $(pwd):/staking_contracts ethereum/solc:0.5.13 \
    @openzeppelin/=/staking_contracts/staking-contract/node_modules/@openzeppelin/ \
    --overwrite -o /staking_contracts/NeuralChainStaking.bin --optimize --optimize-runs 20000 \
    --abi --bin /staking_contracts/staking-contract/contracts/NeuralChainStaking.sol
# Generate go file
echo "Generate go file"
abigen --bin=./NeuralChainStaking.bin/NeuralChainStaking.bin --abi=./NeuralChainStaking.bin/NeuralChainStaking.abi --out ./NeuralChainStaking.go --pkg=staking_contracts
# Generate storage layout file
echo "Generate storage layout file"
cat solc-input.json | docker run -i  -v $(pwd):/staking_contracts -w /staking_contracts ethereum/solc:0.5.13 --standard-json --allow-paths *, > storage-layout.json
popd
