package src

import (
	"fmt"

	model "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/src/model"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	blocks := []model.Block{
		*model.CreateBlock(1, []byte("Genesis Block"), []byte("")),
	}

	for _, block := range blocks {
		serializedBlock, err := block.Serialize()
		if err != nil {
			return fmt.Errorf("failed to serialize block: %v", err)
		}
		err = ctx.GetStub().PutState(string(block.Hash), serializedBlock)
		if err != nil {
			return fmt.Errorf("failed to put block in ledger: %v", err)
		}
	}
	return nil
}

func (s *SmartContract) CreateBlock(ctx contractapi.TransactionContextInterface, blockData string) error {
	block := model.CreateBlock(1, []byte(blockData), []byte("previousHash"))
	existingBlock, err := ctx.GetStub().GetState(string(block.Hash))
	if err != nil {
		return fmt.Errorf("failed to get block: %v", err)
	}
	if existingBlock != nil {
		return fmt.Errorf("block already exists")
	}
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
