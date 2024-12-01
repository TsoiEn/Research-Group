package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/TsoiEn/Research-Group/Soft_Eng-Research/Blockchain_Core/chaincode/consensus" // Replace with the actual import path of your `consensus` package
	//"github.com/TsoiEn/softEng/Blockchain_Core/chaincode/src"
	"github.com/TsoiEn/Research-Group/Soft_Eng-Research/Blockchain_Core/chaincode/src/model" // Replace with the actual import path of your `model` package
)

var studentChain = &model.StudentChain{}

// Function to simulate user input for testing admin operations
// Create a credential blockchain
var credentialChain = &model.CredentialChain{BlockChain: *model.NewBlockChain()}

func testAdminOperations() {
	// Initialize student chain
	if studentChain.Students == nil {
		studentChain.Students = make(map[int]*model.Student)
	}
	// Add a sample student manually to test (if not done elsewhere)
	if _, exists := studentChain.Students[202013432]; !exists {
		studentChain.Students[202013433] = &model.Student{
			StudentID: 202013433,
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	}

	admin := &model.Admin{
		AdminID: "1",
		Name:    "Admin User",
	}

	// Simulate adding a new student
	// Simulate adding a new student
	fmt.Println("Testing AddNewStudent...")
	newStudent, err := admin.AddNewStudent(202013432, "Mark", "Renee", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), 1, studentChain)
	if err == nil && newStudent != nil {
		studentChain.Students[newStudent.StudentID] = newStudent
		fmt.Println("AddNewStudent passed:", newStudent)
	} else {
		fmt.Println("AddNewStudent failed. Error:", err)
	}

	// Display blockchain state after adding student
	displayBlockchainState(&credentialChain.BlockChain)

	// Simulate adding a credential by admin
	fmt.Println("\nTesting AddCredentialAdmin...")
	cred := model.Credential{
		Type:       model.Academic,
		Issuer:     "Admin University",
		DateIssued: time.Now(),
	}
	adminSuccess := admin.AddCredentialAdmin(newStudent, cred.Type, cred.Issuer, cred.DateIssued)
	if adminSuccess {
		// Add credential to the blockchain
		err := credentialChain.AddCredentialModel(&cred)
		if err != nil {
			fmt.Println("Failed to add credential to blockchain:", err)
		} else {
			fmt.Println("AddCredentialAdmin passed.")
		}
	} else {
		fmt.Println("AddCredentialAdmin failed.")
	}

	// Display blockchain state after adding credential
	displayBlockchainState(&credentialChain.BlockChain)
}

func testStudentOperations() {
	// Initialize student chain
	if studentChain.Students == nil {
		studentChain.Students = make(map[int]*model.Student)
	}
	// Add a sample student manually to test (if not done elsewhere)
	if _, exists := studentChain.Students[202013432]; !exists {
		studentChain.Students[202013433] = &model.Student{
			StudentID: 202013433,
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	}

	student := &model.Student{
		StudentID:   202013433,
		FirstName:   "John",
		LastName:    "Doe",
		BirthDate:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Credentials: []*model.Credential{},
	}

	// Simulate adding a credential by a student
	fmt.Println("\nTesting AddCredential...")
	studentID := "202013433"
	cred := model.Credential{
		Type:       model.NonAcademic,
		Issuer:     "Certification Institute",
		DateIssued: time.Now(),
	}
	studentSuccess := student.AddCredential(cred.Type, cred.Issuer, cred.DateIssued)
	if studentSuccess {
		// Add credential to the blockchain
		err := credentialChain.AddCredentialModel(&cred)
		if err != nil {
			fmt.Println("Failed to add credential to blockchain:", err)
		} else {
			fmt.Println("AddCredential passed.")
		}
	} else {
		fmt.Println("AddCredential failed.")
	}

	// Display blockchain state after adding credential
	displayBlockchainState(&credentialChain.BlockChain)

	// Simulate updating student credentials
	fmt.Println("\nTesting UpdateStudentCredentials...")
	studentIDInt, err := strconv.Atoi(studentID)
	if err != nil {
		fmt.Println("Invalid student ID:", studentID)
		return
	}
	updatedSuccess := studentChain.UpdateStudentCredentials(studentIDInt, cred)
	if updatedSuccess {
		fmt.Println("UpdateStudentCredentials passed.")
	} else {
		fmt.Println("UpdateStudentCredentials failed.")
	}

	// Display blockchain state after updating credentials
	displayBlockchainState(&credentialChain.BlockChain)

	// Simulate finding a student by ID
	fmt.Println("\nTesting FindStudentByID...")
	StudentID := studentID
	studentIDInt, err = strconv.Atoi(StudentID)
	if err != nil {
		fmt.Println("Invalid student ID:", StudentID)
		return
	}
	foundStudent, err := studentChain.FindStudentByID(studentIDInt)
	if err != nil {
		fmt.Println("FindStudentByID failed:", err)
	} else {
		fmt.Printf("FindStudentByID passed. Found student: %v\n", foundStudent)
	}
}

// Helper function to display the current state of the blockchain
func displayBlockchainState(blockchain *model.BlockChain) {
	fmt.Println("\n--- Blockchain State ---")
	for i, block := range blockchain.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("  Index: %d\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp)
		fmt.Printf("  Hash: %x\n", block.Hash)
		fmt.Printf("  PrevHash: %x\n", block.PrevHash)
		fmt.Printf("  Data: %s\n\n", string(block.Data))
	}
}

func main() {
	// test consensus
	node1 := consensus.NewRaftNode("node1", []string{"node2", "node3"})
	err := node1.Start()
	if err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	// Manually set node1 as the leader (for testing purposes)
	node1.State = consensus.Leader
	fmt.Println("Node1 is now the leader.")

	// Now you can propose the genesis block
	genesisBlock := &model.Block{
		Index:     0,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      []byte("Genesis Block"),
		PrevHash:  nil,
	}

	success := node1.ProposeBlock(genesisBlock)
	if !success {
		fmt.Println("Error proposing genesis block")
	} else {
		fmt.Println("Genesis block proposed successfully.")
	}

	fmt.Println("\nRunning tests for admin and student operations...")

	// Run admin operations tests
	testAdminOperations()

	// Run student operations tests
	testStudentOperations()

	fmt.Println("\nTesting completed.")

	// test purposes

	// nodeID := "node1"
	// peers := []string{"node2", "node3"}
	// blockchain := src.NewBlockchain(nodeID, peers)

	// // Initialize the ledger with a genesis block.
	// err = blockchain.InitLedger()
	// if err != nil {
	// 	panic(err)
	// }

	// // Example of adding a block.
	// err = blockchain.CreateBlock("First block data")
	// if err != nil {
	// 	panic(err)
	// }

	// peers := []string{"Node1", "Node2", "Node3"}
	// node := consensus.NewRaftNode("Node1", peers)

	// fmt.Println("Starting Raft node...")
	// node.ResetElectionTimer()

	// select {}
}
