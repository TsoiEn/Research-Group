package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Credential struct {
	Type       string    // Type of the credential (e.g., degree, certificate)
	Issuer     string    // Issuer of the credential (e.g., university, organization)
	DataIssued time.Time // Date when the credential was issued
	Hash       []byte    // Hash of the credential data for verification
}

// Serialize converts the Credential to a byte array
func (cred Credential) Serialize() []byte {
	// Convert the Credential to a byte array
	return []byte(fmt.Sprintf("%s|%s|%s", cred.Type, cred.Issuer, cred.DataIssued.Format(time.RFC3339)))
}

// VerifyCredential checks if a credential is valid by comparing its hash
func VerifyCredential(s *Student, cred Credential) bool {
	for _, storedCred := range s.Credentials {
		if bytes.Equal(storedCred.Hash, cred.Hash) {
			return true
		}
	}
	return false
}

func GenerateCredentialHash(cred Credential) []byte {
	// Combine the credential fields to generate a unique hash
	credData := fmt.Sprintf("%s%s%s", cred.Type, cred.Issuer, cred.DataIssued.String())
	hash := sha256.Sum256([]byte(credData)) // Using SHA-256 to hash the credential data
	return hash[:]
}

func ValidateCredentialData(cred Credential) error {
	// Add validation logic for credential data if needed
	// in this case, we can check if the credential type is valid
	if cred.Type == "" {
		return fmt.Errorf("credential type cannot be empty")
	}

	// Check if Issuer is empty
	if cred.Issuer == "" {
		return fmt.Errorf("issuer cannot be empty")
	}

	// Check if DataIssued is a future date
	if cred.DataIssued.After(time.Now()) {
		return fmt.Errorf("issued date cannot be in the future")
	}

	// Additional validation checks can be added here as needed

	return nil
}

//TODO Credential-Related Operations - Any other logic specifically tied to how a credential is created, validated, or manipulated should be moved to credential.go.
