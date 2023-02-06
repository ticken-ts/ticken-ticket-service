package public_blockchain

import (
	"encoding/hex"
	"math/rand"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type DevPublicBlockchain struct {
	contracts map[string]*DevContractCaller
}

func NewDevPublicBlockchain() *DevPublicBlockchain {
	return &DevPublicBlockchain{}
}

// Connect connect to backend
func (pb *DevPublicBlockchain) Connect() error {
	return nil
}

// DeployContract Deploy contract using Geth generated bindings, returns contract address
func (pb *DevPublicBlockchain) DeployContract() (string, error) {
	// Generate random address
	addr, err := randomHex(20)
	if err != nil {
		return "", err
	}

	// Create dev contract instance
	contract, err := NewDevContractCaller(addr)
	if err != nil {
		pb.contracts[addr] = contract
	}

	return addr, nil
}

// GetContract Get contract instance from contract address
func (pb *DevPublicBlockchain) GetContract(contractAddress string) (BCContractCaller, error) {
	return pb.contracts[contractAddress], nil
}

// GeneratePrivateKey Generate private key
func (pb *DevPublicBlockchain) GeneratePrivateKey() (string, error) {
	pk, _ := randomHex(32)
	return pk, nil
}
