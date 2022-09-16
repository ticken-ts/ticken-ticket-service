package pvtbc

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type chaincodeConnector struct {
	chaincode *client.Contract
}

func NewChaincodeConnector(hfc HyperledgerFabricConnector, channelName string, chaincodeName string) (*chaincodeConnector, error) {
	if !hfc.IsConnected() {
		err := hfc.Connect()
		if err != nil {
			return nil, err
		}
	}

	chaincode, err := hfc.GetChaincode(channelName, chaincodeName)
	if err != nil {
		return nil, err
	}

	chaincodeConnector := new(chaincodeConnector)
	chaincodeConnector.chaincode = chaincode

	return chaincodeConnector, nil
}

func (cc *chaincodeConnector) Query(function string, args ...string) ([]byte, error) {
	evaluateResult, err := cc.chaincode.EvaluateTransaction(function, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}

func (cc *chaincodeConnector) Submit(function string, args ...string) ([]byte, error) {
	evaluateResult, err := cc.chaincode.SubmitTransaction(function, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}

func (cc *chaincodeConnector) SubmitAsync(function string, args ...string) ([]byte, *client.Commit) {
	submitResult, commit, err := cc.chaincode.SubmitAsync(function, client.WithArguments(args...))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	return submitResult, commit
}
