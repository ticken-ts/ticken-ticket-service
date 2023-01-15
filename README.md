# Ticken - Ticket Service

## Architectural design

* App > API >

## Smart contract

The Go bindings for the NFT smart contract are committed, but if any changes to the smart contract are made, these need to be updated.

### Requirements

For building the smart contract and regenerating the go bindings.

- Solidity compiler "solc" present in PATH
- Abigen (go-ethereum) "https://geth.ethereum.org/docs/getting-started/installing-geth"

### Build smart contract and go bindings

To compile the smart contract, go to the `infra/public_blockchain/contract` directory and run `npm install`.

Run `npm run build` to compile the smart contract.

Run `npm run abigen` to generate the Go bindings for the smart contract.

## Running locally

This is project is built in way that it can be run locally. 
To achieve this run locally three services:

* Mongo DB instance
* This server (ticken-ticket-service)
* Hyperledger Fabric peer with two chaincodes:
  * ticken-ticket-chaincode
  * ticken-event-chaincode

Before starting clone the following repos in the same folder:
* [ticken-dev](https://github.com/tpp-facu-javi/ticken-dev): contains
all docker images that we are going to use and the scripts to run them.

* [ticken-chaincodes](https://github.com/tpp-facu-javi/ticken-chaincodes): contains 
ticken-event chaincode and ticken-ticket chaincode

All scripts are going to be inside the folder `dev-services` inside `ticken-dev`

### Running the MongoDB instance

```
sh ./start-mongo.sh
```

This is going to start a docker container with a mongo db image.
The image name is `ticken-mongo`

### Running the Hyperledger Fabric Peer

```
sh ./start-pvtbc.sh
```

This is going to start all the images needed to run an hyperledger fabric peer and it
will deploy all necessary chaincodes.

### Running local public blockchain

```bash
sh ./start-public-blockchain.sh
```

### Running ticken-ticket-service

Once you run successfully the hyperledger fabric peer and the MongoDB instance, 
you can start this service

**Running without Docker**

To run this server without docker install the following dependencies:

**Running Docker**

## Running tests

**Running especific package**

Use the following commnad to run the test in specific package
```
go test ./<paht_to_package>
```

**Running all tests**

Use the following commnad to run all tests in the project
```
go test ./...
```