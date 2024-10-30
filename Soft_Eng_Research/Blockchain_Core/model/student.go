package blockchain

import (
	"bytes"
	"fmt"
	"time"
)

type Student struct {
	ID          int
	FirstName   string
	LastName    string
	Age         int
	BirthDate   string
	StudentNum  int
	Credentials []Credential
}

// AddNewStudent creates and returns a new student
func AddNewStudent(id int, firstName, lastName string, age int, birthDate string, studentNum int) *Student {
	student := &Student{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		Age:         age,
		BirthDate:   birthDate,
		StudentNum:  studentNum,
		Credentials: []Credential{},
	}
	return student
}

// AddCredential adds a credential to the student and generates a hash
// specifically non-academic
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
		return err // Return the validation error
	}

	// Generate and store the credential hash
	newCredential.Hash = GenerateCredentialHash(newCredential)

	// Add the credential to the student's list of credentials
	s.Credentials = append(s.Credentials, newCredential)
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
	student.Credentials = append(student.Credentials, newCredential)
	// Return successfully
	return nil
}

// FindStudentByID should find and return a student by ID (this needs to be implemented)
func FindStudentByID(id int) (*Student, error) {
	// Logic to find the student by ID from the blockchain
	// For now, return nil and an error to avoid compilation issues
	Students := []*Student{}
	for _, student := range Students {
		if student.ID == id {
			return student, nil
		}
	}
	return nil, fmt.Errorf("Student not found")
}

// DeleteCredential removes a credential from the student's list
func DeleteCredential(s *Student, cred Credential) {
	for i, storedCred := range s.Credentials {
		// Correct comparison using bytes.Equal
		if bytes.Equal(storedCred.Hash, cred.Hash) {
			s.Credentials = append(s.Credentials[:i], s.Credentials[i+1:]...)
			break
		}
	}
}

// LookupCredentials retrieves and displays all the credentials of the student
func LookupCredentials(s *Student) {
	for _, cred := range s.Credentials {
		fmt.Printf("Credential: %s, Issued by: %s, Issued on: %s\n", cred.Type, cred.Issuer, cred.DateIssued.String())
	}
}
