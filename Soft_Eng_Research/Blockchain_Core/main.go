package main

import (
	"fmt"
	"log"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MockTransactionContext struct {
	contractapi.TransactionContextInterface
}

func (m *MockTransactionContext) GetStub() *contractapi.TransactionContextInterface {
	// Return a mock stub if needed
	return nil
}

func main() {
	// Create a new chaincode instance
	cc := new(chaincode.Chaincode)

	// Mock student login
	var loginID string
	var loginPass string

	// Mock UI login section
	fmt.Print("Student Number: ")
	fmt.Scanln(&loginID)
	fmt.Print("Password: ")
	fmt.Scanln(&loginPass)

	// Mock student login data
	students := map[string]string{
		"202533282": "zpKVM4cQ",
		"202403450": "S2oS6gQP",
		"202209675": "WXwoXDA9",
		"202433194": "uOlgXCpt",
		"202226488": "TOTVufqI",
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

	// Create a mock transaction context
	mockCtx := &MockTransactionContext{}

	// Simulated credential data
	credentialData := "Degree in Computer Science"

	// Mock updating a credential
	if err := cc.UpdateCredential(mockCtx, loginID, credentialData); err != nil {
		log.Fatalf("Failed to update credential: %v", err)
	}

	// Mock verifying a credential
	isValid, err := cc.VerifyCredential(mockCtx, loginID, credentialData)
	if err != nil {
		log.Fatalf("Failed to verify credential: %v", err)
	}
	if isValid {
		fmt.Println("Credential is valid.")
	} else {
		fmt.Println("Credential is not valid.")
	}

	// Mock retrieving credentials
	credentials, err := cc.RetrieveCredential(mockCtx, loginID)
	if err != nil {
		log.Fatalf("Failed to retrieve credentials: %v", err)
	}

	// Print out the retrieved credentials
	fmt.Printf("Credentials for student %s:\n", loginID)
	for _, cred := range credentials {
		fmt.Println(cred)
	}
}
