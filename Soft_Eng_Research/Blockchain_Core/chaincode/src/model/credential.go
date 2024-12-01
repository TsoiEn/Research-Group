package model

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

// CredentialType enumeration
type CredentialType int

const (
	Academic CredentialType = iota
	NonAcademic
	Certificate
	Diploma
)

func (ct CredentialType) String() string {
	return [...]string{"Academic", "NonAcademic", "Certificate", "Diploma"}[ct]
}

// Credential represents an individual credential.
type Credential struct {
	ID         string         `json:"id"`
	Type       CredentialType `json:"type"`
	Issuer     string         `json:"issuer"`
	DateIssued time.Time      `json:"date_issued"`
	Hash       []byte         `json:"hash"`
	Status     string         `json:"status"`
}

// ValidateCredentialData ensures the credential fields are valid.
func ValidateCredentialData(cred *Credential) error {
	if cred.Type.String() == "" {
		return fmt.Errorf("credential type cannot be empty")
	}
	if cred.Issuer == "" {
		return fmt.Errorf("issuer cannot be empty")
	}
	if cred.DateIssued.After(time.Now()) {
		return fmt.Errorf("issued date cannot be in the future")
	}
	return nil
}

// CredentialChain is an alias for BlockChain, which stores credentials.
type CredentialChain struct {
	BlockChain
}

// AddCredential adds a new credential to the blockchain.
// func (chain *CredentialChain) AddCredentialToBlockchain(cred *Credential) error {
// 	if err := ValidateCredentialData(cred); err != nil {
// 		return err
// 	}
// 	cred.Hash = GenerateCredentialHash(cred)
// 	credData, err := json.Marshal(cred)
// 	if err != nil {
// 		return err
// 	}
// 	chain.AddBlock(credData)
// 	return nil
// }

// AddCredentialToBlockchain simulates adding a credential to the blockchain.
func AddCredentialToBlockchain(ownerID string, fileBytes []byte, filetype, credentialType, issuer, dateIssued string) {
	// Log the credential information being added to the blockchain
	log.Printf("Adding credential to blockchain: ownerID=%s, filetype=%s, credentialType=%s, issuer=%s, dateIssued=%s", ownerID, filetype, credentialType, issuer, dateIssued)

	// Convert dateIssued to time.Time
	issuedDate, err := time.Parse("2006-01-02", dateIssued)
	if err != nil {
		log.Fatalf("Invalid date format: %s, expected YYYY-MM-DD", dateIssued)
		return
	}

	// Create a new Credential
	cred := &Credential{
		ID:         ownerID, // Assuming ownerID is unique for each credential
		Type:       mapCredentialType(credentialType),
		Issuer:     issuer,
		DateIssued: issuedDate,
		Hash:       fileBytes, // Assuming fileBytes is the credential's hash or relevant data
		Status:     "Active",  // Default status, could be more dynamic based on your use case
	}

	// Add the credential to the blockchain (simulating here)
	if err := AddCredentialToBlockchainLogic(cred); err != nil {
		log.Printf("Error adding credential: %v", err)
	}
}

// mapCredentialType maps the string value to CredentialType enum.
func mapCredentialType(credentialType string) CredentialType {
	switch credentialType {
	case "Academic":
		return Academic
	case "NonAcademic":
		return NonAcademic
	case "Certificate":
		return Certificate
	case "Diploma":
		return Diploma
	default:
		log.Printf("Unknown credential type: %s", credentialType)
		return Academic // Default to Academic if unknown
	}
}

// AddCredentialToBlockchainLogic adds the credential to the blockchain (simulated).
func AddCredentialToBlockchainLogic(cred *Credential) error {
	// This is a placeholder function to simulate adding to the blockchain
	log.Printf("Credential added to blockchain: %+v", cred)
	return nil
}

// VerifyCredential checks if a credential exists in the blockchain.
func (chain *CredentialChain) VerifyCredential(id string) (bool, error) {
	cred, err := chain.FindCredentialByID(id)
	if err != nil {
		return false, err
	}
	expectedHash := GenerateCredentialHash(cred)
	return bytes.Equal(cred.Hash, expectedHash), nil
}
