package model

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"time"
)

type CredentialType int

// CredentialType enumeration
const (
	Academic CredentialType = iota
	NonAcademic
	Certificate
	Diploma
)

// String returns the string representation of the CredentialType
func (ct CredentialType) String() string {
	return [...]string{"Academic", "NonAcademic", "Certificate", "Diploma"}[ct]
}

type Credential struct {
	ID         string         // Unique identifier for the credential
	Type       CredentialType // Type of the credential (Academic or NonAcademic)
	Issuer     string         // Issuer of the credential (e.g., university, organization)
	DateIssued time.Time      // Date when the credential was issued
	Hash       []byte         // Hash of the credential data for verification
	PrevHash   []byte         // Hash of the previous credential
}

// CredentialChain represents a chain of credentials
type CredentialChain struct {
	Credentials []Credential
}

// Serialize converts the Credential to a custom byte array format
func (cred Credential) Serialize() []byte {
	return []byte(fmt.Sprintf("%d|%s|%s|%s", cred.Type, cred.Issuer, cred.ID, cred.DateIssued.Format(time.RFC3339)))
}

// VerifyCredential checks if a credential is valid by comparing its hash
func (cred *Credential) VerifyCredential(s *Student) bool {
	expectedHash := GenerateCredentialHash(*cred)
	for _, storedCred := range s.Credentials {
		if bytes.Equal(storedCred.Hash, expectedHash) {
			return true
		}
	}
	log.Println("Credential verification failed")
	return false
}

// GenerateCredentialHash creates a hash of the credential data for integrity
func GenerateCredentialHash(cred Credential) []byte {
	credData := fmt.Sprintf("%d|%s|%s|%s", cred.Type, cred.Issuer, cred.ID, cred.DateIssued.String())
	hash := sha256.Sum256([]byte(credData))
	return hash[:]
}

// ValidateCredentialData checks the validity of the credential fields
func ValidateCredentialData(cred Credential) error {
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
