package public_blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	public_blockchain "ticken-ticket-service/infra/public_blockchain/contract"
	"ticken-ticket-service/models"
)

type ContractCaller struct {
	addr       string
	instance   *public_blockchain.TickenEvent
	conn       Connection
	transactor *bind.TransactOpts
}

// NewContractCaller creates a new contract caller instance with the given contract address
func NewContractCaller(contractAddr string, conn Connection, transactor *bind.TransactOpts) (*ContractCaller, error) {
	instance, err := public_blockchain.NewTickenEvent(common.HexToAddress(contractAddr), conn)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{
		instance:   instance,
		conn:       conn,
		transactor: transactor,
		addr:       contractAddr,
	}, nil
}

// GenerateTicket BLOCKING generate ticket and assign the buyer as owner
func (cc *ContractCaller) GenerateTicket(
	buyerAddress string,
	ticketData *models.Ticket,
) (*string, error) {
	tx, err := cc.instance.SafeMint(
		cc.transactor,
		common.HexToAddress(buyerAddress),
		ticketData.Section,
	)
	if err != nil {
		return nil, err
	}

	// Wait for transaction to be mined
	_, err = bind.WaitMined(cc.transactor.Context, cc.conn, tx)
	if err != nil {
		return nil, err
	}

	transactionAddress := tx.Hash().String()
	return &transactionAddress, nil
}

// GetUserTickets returns all tickets owned by the given user
func (cc *ContractCaller) GetUserTickets(userAddress string) ([]PubBCTicket, error) {
	tickets, err := cc.instance.GetTicketsByOwner(nil, common.HexToAddress(userAddress))
	if err != nil {
		return nil, err
	}

	var result []PubBCTicket
	for _, ticket := range tickets {
		result = append(result, PubBCTicket{
			Section:         ticket.Section,
			TokenID:         ticket.TokenID.String(),
			OwnerAddress:    userAddress,
			ContractAddress: cc.addr,
		})
	}

	return result, nil
}
