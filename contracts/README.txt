Run the `generated_docker.sh` script to generate Golang files from the smart contracts.

Ensure that you have the Perun smart contracts placed in the `perun-eth-contracts` directory within the `contracts` folder. ([`Perun Contracts Github`](https://github.com/hyperledger-labs/perun-eth-contracts))

Additionally, make sure that `perun-eth-contracts` includes `openzeppelin-contracts` **(version < 5.0)**. ([`Openzeppelin Contracts Github`](https://github.com/OpenZeppelin/openzeppelin-contracts))

Ensure you have the docker image `ethereum/solc:0.8.24` to run the `generate_docker.sh` script.
([`How to install ethereum/solc using docker`](https://docs.soliditylang.org/en/latest/installing-solidity.html#docker))