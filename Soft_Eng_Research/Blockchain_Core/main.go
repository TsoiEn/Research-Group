package main

import (
	"fmt"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/blockchain"
)

func main() {
	// Create a new blockchain with the genesis block
	chain := blockchain.NewBlockChain()

	// Add new blocks to the blockchain
	chain.AddBlock("Choi a member of this group")
	chain.AddBlock("Raffy a member of this group")
	chain.AddBlock("Yonne a member of this group")

	// Print out the details of each block in the blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in the block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
