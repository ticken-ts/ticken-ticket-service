package public_blockchain

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

type Connection interface {
	bind.ContractBackend
	bind.DeployBackend
}
