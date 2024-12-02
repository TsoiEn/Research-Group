package src

import (
	"fmt"
	"log"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/consensus"
	model "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/src/model"
)

type Blockchain struct {
	Blocks   []*model.Block      // Blockchain blocks
	RaftNode *consensus.RaftNode // Raft node for consensus
}

// Initialize the blockchain with a genesis block.
func (bc *Blockchain) InitLedger() error {
	// Create the genesis block.
	genesisBlock := model.CreateBlock(1, []byte("Genesis Block"), []byte(""))

	// Use Raft to propose the genesis block.
	if bc.RaftNode == nil {
		return fmt.Errorf("RaftNode is not initialized")
	}

	success := bc.RaftNode.ProposeBlock(genesisBlock)
	if !success {
		return fmt.Errorf("failed to propose genesis block via Raft")
	}

	log.Printf("Genesis block initialized: %+v", genesisBlock)
	return nil
}

// Create a new block and add it to the blockchain.
func (bc *Blockchain) CreateBlock(data string) error {
	// Ensure the blockchain is initialized.
	if len(bc.Blocks) == 0 {
		return fmt.Errorf("blockchain is not initialized")
	}
	// Step 1: Get the last block.
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	// Step 2: Create a new block using the previous block's hash.
	newBlock := model.CreateBlock(lastBlock.Index+1, []byte(data), lastBlock.Hash)

	// Step 3: Propose the new block via Raft.
	success := bc.RaftNode.ProposeBlock(newBlock)
	if !success {
		return fmt.Errorf("failed to propose block via Raft")
	}

	// Step 4: Append the new block to the blockchain after consensus.
	bc.Blocks = append(bc.Blocks, newBlock)
	log.Printf("New block added via Raft: %+v", newBlock)
	return nil
}

func NewBlockchain(nodeID string, peers []string) *Blockchain {
	// Step 1: Create and start the Raft node
	raftNode := consensus.NewRaftNode(nodeID, peers)

	// Start the Raft node
	err := raftNode.Start()
	if err != nil {
		log.Fatalf("Failed to start Raft node: %v", err)
	}

	genesisBlock := model.CreateBlock(0, []byte("Genesis Block"), nil)

	chain := &Blockchain{
		Blocks:   []*model.Block{genesisBlock},
		RaftNode: raftNode,
	}

	log.Printf("Blockchain initialized with genesis block: %+v", genesisBlock)

	return chain
}
