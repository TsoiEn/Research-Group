package main

import (
	"fmt"

	blockchain "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/model"
)

func main() {
	// Create a new blockchain with the genesis block
	chain := blockchain.NewBlockChain()

	// mock student log in
	loginID := 0
	loginPass := ""

	// mock ui log in section
	fmt.Print("Student Number: ")
	fmt.Scanln(&loginID)
	fmt.Print("Password: ")
	fmt.Scanln(&loginPass)

	// Mock student login data
	students := map[int]string{
		202533282: "zpKVM4cQ",
		202403450: "S2oS6gQP",
		202209675: "WXwoXDA9",
		202433194: "uOlgXCpt",
		202226488: "TOTVufqI",
	}

	// Check if the login ID exists and the password matches
	if pass, exists := students[loginID]; exists {
		if pass == loginPass {
			fmt.Println("Login successful!")
		} else {
			fmt.Println("Invalid password.")
			return
		}
	} else {
		fmt.Println("Student ID not found.")
		return
	}

	// ! commented for now

	// // Add a new student
	// student := blockchain.AddNewStudent(1, "John", "Doe", 21, "Computer Science", 2023)

	// // Add credentials to the student
	// student.AddCredential(blockchain.Degree, "University A", time.Now())
	// student.AddCredential(blockchain.Transcript, "University A", time.Now())

	// Iterate over the student's credentials and add them to the blockchain
	// Serializing each credential
	// Add to blockchain
	// for _, cred := range student.Credentials {
	// 	blockData := cred.Serialize()
	// 	chain.AddBlock(blockData)
	// }

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
