package public_blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	public_blockchain "ticken-ticket-service/infra/public_blockchain/contract"
)

type PublicBlockchain struct {
	chainUrl  string
	chainID   int64
	addressPK string
	auth      *bind.TransactOpts
	conn      Connection
}

func NewPublicBlockchain(chainUrl string, chainID int64, addressPK string) *PublicBlockchain {
	return &PublicBlockchain{
		chainUrl:  chainUrl,
		chainID:   chainID,
		addressPK: addressPK,
	}
}

// Connect connect to backend
func (pb *PublicBlockchain) Connect() error {
	auth, err := pb.getTransactor(pb.addressPK)
	if err != nil {
		println(err.Error())
		return err
	}

	conn, err := ethclient.Dial(pb.chainUrl)
	if err != nil {
		println(err.Error())
		return err
	}

	pb.auth = auth
	pb.conn = conn
	return nil
}

// getTransactor convert pk as hex string to a transactor object for contract calls
func (pb *PublicBlockchain) getTransactor(pk string) (*bind.TransactOpts, error) {
	ecdsaKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactorWithChainID(ecdsaKey, big.NewInt(pb.chainID))
}

// DeployContract Deploy contract using Geth generated bindings, returns contract address
func (pb *PublicBlockchain) DeployContract() (*string, error) {
	addr, tx, _, err := public_blockchain.DeployTickenEvent(pb.auth, pb.conn)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	_, err = bind.WaitDeployed(pb.auth.Context, pb.conn, tx)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	addrString := addr.String()

	return &addrString, nil
}

// GetContract Get contract instance from contract address
func (pb *PublicBlockchain) GetContract(contractAddress string) (*ContractCaller, error) {
	return NewContractCaller(contractAddress, pb.conn)
}
