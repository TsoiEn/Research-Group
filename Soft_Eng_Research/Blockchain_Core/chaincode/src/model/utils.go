package model

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Serialize converts the Credential to a custom byte array format
func (cred *Credential) Serialize() []byte {
	return []byte(fmt.Sprintf("%d|%s|%s|%s", cred.Type, cred.Issuer, cred.ID, cred.DateIssued.Format(time.RFC3339)))
}

// GenerateCredentialHash creates a hash of the credential data for integrity
func GenerateCredentialHash(cred *Credential) []byte {
	credData := cred.Serialize()
	hash := sha256.Sum256(credData)
	return hash[:]
}
