package model

import (
	"bytes"
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
	ID         string         `json:"id"`
	Type       CredentialType `json:"type"`
	Issuer     string         `json:"issuer"`
	DateIssued time.Time      `json:"date_issued"`
	Hash       []byte         `json:"hash"`
	PrevHash   []byte         `json:"prev_hash"`
	Status     string         `json:"status"`
}

// CredentialChain represents a chain of credentials
type CredentialChain struct {
	Credentials []Credential
}

// VerifyCredential checks if a credential is valid by comparing its hash
func (cred *Credential) VerifyCredential(s *Student) bool {
	expectedHash := GenerateCredentialHash(cred)
	for _, storedCred := range s.Credentials {
		if bytes.Equal(storedCred.Hash, expectedHash) {
			return true
		}
	}
	log.Println("Credential verification failed")
	return false
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

func LookupCredentials(s *Student) {
	for _, cred := range s.Credentials {
		fmt.Printf("Credential: %s, Issued by: %s, Issued on: %s, Status: %s\n",
			cred.Type, cred.Issuer, cred.DateIssued.String(), cred.Status)
	}
}
