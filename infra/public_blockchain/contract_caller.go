package public_blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	public_blockchain "ticken-ticket-service/infra/public_blockchain/contract"
)

type ContractCaller struct {
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
func (cc *ContractCaller) GenerateTicket(buyerAddress string, infoUrl string) (*string, error) {
	tx, err := cc.instance.SafeMint(cc.transactor, common.HexToAddress(buyerAddress), infoUrl)
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
