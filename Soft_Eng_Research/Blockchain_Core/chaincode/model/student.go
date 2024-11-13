package model

import (
	"bytes"
	"fmt"
	"time"
)

type Student struct {
	StudentID   int    `json:"student_id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	Age         int
	BirthDate   string        `json:"birth_date"`
	Credentials []*Credential `json:"credentials,omitempty"`
}

// AddCredential adds a new credential to the student's list of non-academic credentials
func (s *Student) AddCredential(credentialType CredentialType, issuer string, dataIssued time.Time) error {
	// Check if the credential type is non-academic
	if credentialType != NonAcademic {
		return fmt.Errorf("only non-academic credentials can be added")
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
	newCredential.Hash = GenerateCredentialHash(&newCredential)

	// Add the credential to the student's list of credentials
	s.Credentials = append(s.Credentials, &newCredential)
	return nil
}

// UpdateStudentCredentials updates the credentials of the student
func UpdateStudentCredentials(id int, newCredential Credential) error {
	student, err := FindStudentByID(id)
	if err != nil {
		return err
	}

	// Check if the credential already exists by comparing hashes
	for _, cred := range student.Credentials {
		if bytes.Equal(cred.Hash, newCredential.Hash) {
			return fmt.Errorf("Credential already exists")
		}
	}

	// Add new credential
	student.Credentials = append(student.Credentials, &newCredential)
	// Return successfully
	return nil
}

// FindStudentByID should find and return a student by ID (this needs to be implemented)
func FindStudentByID(id int) (*Student, error) {
	// Logic to find the student by ID from the blockchain
	// For now, return nil and an error to avoid compilation issues
	Students := []*Student{}
	for _, student := range Students {
		if student.StudentID == id {
			return student, nil
		}
	}
	return nil, fmt.Errorf("Student not found")
}

// DeleteCredential removes a credential from the student's list
