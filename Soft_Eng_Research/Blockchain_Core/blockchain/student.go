package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Student struct {
	ID          int
	FirstName   string
	LastName    string
	Age         int
	BirthDate   time.Time
	Credentials []Credential
}

type Credential struct {
	Type       string
	Issuer     string
	DataIssued time.Time
	Hash       []byte
}

// Serialize converts the Credential to a byte array
func (cred Credential) Serialize() []byte {
	// Convert the Credential to a byte array
	return []byte(fmt.Sprintf("%s|%s|%s", cred.Type, cred.Issuer, cred.DataIssued.Format(time.RFC3339)))
}

// AddNewStudent creates and returns a new student
func AddNewStudent(id int, firstName, lastName string, age int, birthDate time.Time) *Student {
	student := &Student{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		Age:         age,
		BirthDate:   birthDate,
		Credentials: []Credential{},
	}
	return student
}

// AddCredential adds a credential to the student and generates a hash
func (s *Student) AddCredential(credentialType, issuer string, dataIssued time.Time) error {
	// Create a new credential
	newCredential := Credential{
		Type:       credentialType,
		Issuer:     issuer,
		DataIssued: dataIssued,
	}

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

// VerifyCredential checks if a credential is valid by comparing its hash
func VerifyCredential(s *Student, cred Credential) bool {
	for _, storedCred := range s.Credentials {
		if bytes.Equal(storedCred.Hash, cred.Hash) {
			return true
		}
	}
	return false
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

// GenerateCredentialHash creates a hash for a credential
func GenerateCredentialHash(cred Credential) []byte {
	// Combine the credential fields to generate a unique hash
	credData := fmt.Sprintf("%s%s%s", cred.Type, cred.Issuer, cred.DataIssued.String())
	hash := sha256.Sum256([]byte(credData)) // Using SHA-256 to hash the credential data
	return hash[:]
}

// ValidateCredentialData validates the data of the credential (to be implemented)
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

// LookupCredentials retrieves and displays all the credentials of the student
func LookupCredentials(s *Student) {
	for _, cred := range s.Credentials {
		fmt.Printf("Credential: %s, Issued by: %s, Issued on: %s\n", cred.Type, cred.Issuer, cred.DataIssued.String())
	}
}
