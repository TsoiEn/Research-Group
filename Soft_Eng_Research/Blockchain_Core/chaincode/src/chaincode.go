package src

import (
	"fmt"
	"log"

	model "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/src/model"
)

type Blockchain struct {
	Blocks []*model.Block
}

// Initialize the blockchain with a genesis block.
func (bc *Blockchain) InitLedger() error {
	// Create the genesis block.
	genesisBlock := model.CreateBlock(1, []byte("Genesis Block"), []byte(""))

	// Append the genesis block to the blockchain.
	bc.Blocks = append(bc.Blocks, genesisBlock)
	log.Printf("Genesis block initialized: %+v", genesisBlock)
	return nil
}

// Create a new block and add it to the blockchain.
func (bc *Blockchain) CreateBlock(data string) error {
	if len(bc.Blocks) == 0 {
		return fmt.Errorf("blockchain is not initialized")
	}

	// Get the last block.
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	// Create a new block using the previous block's hash.
	newBlock := model.CreateBlock(lastBlock.Index+1, []byte(data), lastBlock.Hash)

	// Check for duplicate blocks (using hash).
	for _, block := range bc.Blocks {
		if string(block.Hash) == string(newBlock.Hash) {
			return fmt.Errorf("block already exists with hash: %x", newBlock.Hash)
		}
	}

	// Append the new block to the blockchain.
	bc.Blocks = append(bc.Blocks, newBlock)
	log.Printf("New block created: %+v", newBlock)
	return nil
}
