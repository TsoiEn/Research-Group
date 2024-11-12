package main

import (
	"fmt"

	model "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/model"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Code for ledger initialization
	return nil
}

func (s *SmartContract) CreateBlock(ctx contractapi.TransactionContextInterface, blockData string) error {
	block := model.CreateBlock(1, []byte(blockData), []byte("previousHash"))
	// Code to add block to the ledger
	serializedBlock, err := block.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}
	return ctx.GetStub().PutState(string(block.Hash), serializedBlock)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
