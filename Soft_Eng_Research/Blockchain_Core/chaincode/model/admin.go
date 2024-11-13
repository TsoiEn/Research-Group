package model

import (
	"bytes"
	"fmt"
	"time"
)

type Admin struct {
	AdminID string `json:"admin_id"`
	Name    string `json:"name"`
}

func (a *Admin) AddNewStudent(id int, firstName, lastName string, age int, birthDate string, studentNum int) *Student {
	student := &Student{
		StudentID:   id,
		FirstName:   firstName,
		LastName:    lastName,
		Age:         age,
		BirthDate:   birthDate,
		Credentials: []*Credential{},
	}
	return student
}

// AddCredentialAdmin adds a new academic credential to the student's list of academic credentials
func (a *Admin) AddCredentialAdmin(s *Student, credentialType CredentialType, issuer string, dataIssued time.Time) error {
	// Check if the credential type is academic
	if credentialType != Academic {
		return fmt.Errorf("only academic credentials can be added")
	}

	// Create a new credential
	newCredential := Credential{
		Type:       credentialType,
		Issuer:     issuer,
		DateIssued: dataIssued,
	}

	// Validate the credential data
	if err := ValidateCredentialData(newCredential); err != nil {
		return err
	}

	// Generate and store the credential hash
	newCredential.Hash = GenerateCredentialHash(newCredential)

	// Add the credential to the student's list of credentials
	s.Credentials = append(s.Credentials, &newCredential)
	return nil
}

// RevokeCredential revokes a credential of the student
func RevokeCredential(s *Student, cred Credential) error {
	for _, storedCred := range s.Credentials {
		// Check if the hash matches to identify the credential
		if bytes.Equal(storedCred.Hash, cred.Hash) {
			if storedCred.Status == "revoked" {
				return fmt.Errorf("credential is already revoked")
			}

			// Mark the credential as revoked
			storedCred.Status = "revoked"
			return nil
		}
	}
	return fmt.Errorf("credential not found")
}
