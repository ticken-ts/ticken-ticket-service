package public_blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
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
	}, nil
}

//GenerateTicket BLOCKING generate ticket and assign the buyer as owner
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

// WatchTicketCreatedEvent watch for ticket created event and send it to callback
func (cc *ContractCaller) WatchTicketCreatedEvent(
	callback func(ticketData *CreatedTicket),
) error {
	// Watch for TicketCreated event
	channel := make(chan *public_blockchain.TickenEventTicketCreated)
	_, err := cc.instance.WatchTicketCreated(
		&bind.WatchOpts{},
		channel,
		[]*big.Int{},
	)
	if err != nil {
		return err
	}

	// Wait for event
	go func() {
		for {
			select {
			case event := <-channel:
				ticketData := &CreatedTicket{
					ContractAddr: cc.addr,
					Section:      event.Section,
					Owner:        event.OwnerAddress.String(),
					TokenID:      event.TokenID.Uint64(),
					TxAddress:    event.Raw.TxHash.String(),
				}
				callback(ticketData)
			}
		}
	}()

	return nil
}
