package model

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func (chain *CredentialChain) AddCredential(cred *Credential) error {
	if err := ValidateCredentialData(cred); err != nil {
		return err
	}
	cred.Hash = GenerateCredentialHash(cred)
	credData, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	chain.AddBlock(credData)
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
