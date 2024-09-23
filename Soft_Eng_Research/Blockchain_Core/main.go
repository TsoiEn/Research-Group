package main

import (
	"fmt"
	"time"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/blockchain"
)

func main() {
	// Create a new blockchain with the genesis block
	chain := blockchain.NewBlockChain()

	// Add a new student
	student := blockchain.AddNewStudent(1, "John", "Doe", 21, time.Now())

	// Add credentials to the student
	student.AddCredential("Degree", "University A", time.Now())
	student.AddCredential("Transcript", "University A", time.Now())

	// Iterate over the student's credentials and add them to the blockchain
	for _, cred := range student.Credentials {
		blockData := cred.Serialize()                                        // Serialize each credential
		chain.AddBlock("Added credential for "+student.FirstName, blockData) // Add to blockchain
	}

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
