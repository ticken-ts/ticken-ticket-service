package public_blockchain

import (
	"crypto/rand"
	"math/big"
	"ticken-ticket-service/models"
)

type DevContractCaller struct {
	addr    string
	tickets map[*big.Int]PubBCTicket
}

// NewDevContractCaller creates a new contract caller instance with the given contract address
func NewDevContractCaller(contractAddr string) (*DevContractCaller, error) {
	return &DevContractCaller{
		addr: contractAddr,
	}, nil
}

// GenerateTicket BLOCKING generate ticket and assign the buyer as owner
func (cc *DevContractCaller) GenerateTicket(
	buyerAddress string,
	ticketData *models.Ticket,
) (*string, error) {
	// Generate random token ID
	tokenID, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
	if err != nil {
		return nil, err
	}

	// Generate mocked ticket
	ticket := PubBCTicket{
		Section:         ticketData.Section,
		OwnerAddress:    buyerAddress,
		TokenID:         tokenID,
		ContractAddress: cc.addr,
	}

	// Add ticket to the list
	cc.tickets[tokenID] = ticket

	// Generate random transaction hash
	txHash, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
	if err != nil {
		return nil, err
	}

	transactionAddress := txHash.String()
	return &transactionAddress, nil
}

// GetUserTickets returns all tickets owned by the given user
func (cc *DevContractCaller) GetUserTickets(userAddress string) ([]PubBCTicket, error) {
	var tickets []PubBCTicket
	for _, ticket := range cc.tickets {
		if ticket.OwnerAddress == userAddress {
			tickets = append(tickets, ticket)
		}
	}
	return tickets, nil
}
