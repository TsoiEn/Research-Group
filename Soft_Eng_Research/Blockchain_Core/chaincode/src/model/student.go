package model

import (
	"bytes"
	"fmt"
	"time"
)

type Student struct {
	ID          int           `json:"user_id"`
	LastName    string        `json:"last_name"`
	FirstName   string        `json:"first_name"`
	BirthDate   time.Time     `json:"birth_date"`
	StudentID   int           `json:"student_id"`
	Credentials []*Credential `json:"credentials,omitempty"`
}

type StudentChain struct {
	Students map[int]*Student
}

// AddCredential adds a new credential to the student's list of non-academic credentials
func (s *Student) AddCredential(credentialType CredentialType, issuer string, dataIssued time.Time) bool {
	// Check if the credential type is valid
	if credentialType != NonAcademic && credentialType != Academic {
		return false //fmt.Errorf("invalid credential type")
	}

	// Check if the credential type is non-academic
	if credentialType != NonAcademic {
		return false //fmt.Errorf("only non-academic credentials can be added")
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

// UpdateStudentCredentials updates the credentials of the student
func (chain *StudentChain) UpdateStudentCredentials(id int, newCredential Credential) bool {
	student, err := chain.FindStudentByID(id)
	if err != nil {
		return false // Student not found
	}

	// Check if the credential already exists by comparing hashes
	for _, cred := range student.Credentials {
		if bytes.Equal(cred.Hash, newCredential.Hash) && cred.Type == newCredential.Type && cred.Issuer == newCredential.Issuer && cred.DateIssued.Equal(newCredential.DateIssued) {
			return false //fmt.Errorf("Credential already exists")
		}
	}

	// Add new credential
	student.Credentials = append(student.Credentials, &newCredential)
	// Return successfully
	return true

}

// FindStudentByID should find and return a student by ID (this needs to be implemented)
func (chain *StudentChain) FindStudentByID(id int) (*Student, error) {
	for _, student := range chain.Students {
		if student.StudentID == id {
			return student, nil
		}
	}
	return nil, fmt.Errorf("Student not found")
}
