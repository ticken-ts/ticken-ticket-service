package public_blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"ticken-ticket-service/models"
)

type Connection interface {
	bind.ContractBackend
	bind.DeployBackend
}

type BCContractCaller interface {
	GenerateTicket(buyerAddress string, ticketData *models.Ticket) (*string, error)
	GetUserTickets(userAddress string) ([]PubBCTicket, error)
}

type PublicBC interface {
	Connect() error
	DeployContract() (string, error)
	GetContract(contractAddress string) (BCContractCaller, error)
	GeneratePrivateKey() (string, error)
	GetAddressFromPK(pk string) (string, error)
}
