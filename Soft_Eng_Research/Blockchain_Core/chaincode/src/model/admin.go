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

func (a *Admin) AddNewStudent(id int, firstName, lastName string, birthDate time.Time, studentNum int, chain *StudentChain) *Student {
	if chain.Students == nil {
		chain.Students = make(map[int]*Student)
	}
	student := Student{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		StudentID: studentNum,
	}
	chain.Students[id] = &student
	return &student
}

// AddCredentialAdmin adds a new academic credential to the student's list of academic credentials
func (a *Admin) AddCredentialAdmin(s *Student, credentialType CredentialType, issuer string, dataIssued time.Time) bool {
	// Check if the credential type is academic
	if credentialType != Academic {
		return false //fmt.Errorf("only academic credentials can be added")
	}

	// Create a new credential
	newCredential := Credential{
		Type:       credentialType,
		Issuer:     issuer,
		DateIssued: dataIssued,
	}

	// Validate the credential data
	if err := ValidateCredentialData(&newCredential); err != nil {
		return false
	}

	// Generate and store the credential hash
	newCredential.Hash = GenerateCredentialHash(&newCredential)

	// Add the credential to the student's list of credentials
	s.Credentials = append(s.Credentials, &newCredential)
	return true
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
