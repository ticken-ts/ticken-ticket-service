package public_blockchain

import "math/big"

type PubBCTicket struct {
	Section         string
	OwnerAddress    string
	TokenID         *big.Int
	ContractAddress string
}
