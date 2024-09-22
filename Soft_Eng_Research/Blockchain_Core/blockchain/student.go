package blockchain

import (
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

func (s *Student) AddCredential(credentialType, issuer string, dataIssued time.Time, hash []byte) {
	newCredential := Credential{
		Type:       credentialType,
		Issuer:     issuer,
		DataIssued: dataIssued,
		Hash:       hash,
	}
	s.Credentials = append(s.Credentials, newCredential)
}

func UpdateStudentCredentials(id int, newCredential Credential) error {
	student, err := FindStudentByID(id)
	if err != nil {
		return err
	}

	for _, cred := range student.Credentials {
		if cred.Type == newCredential.Hash {
			return fmt.Errorf("Credential already exists")
		}
	}

	student.Credentials = append(student.Credentials, newCredential)
	//! return studentCredentials
}

func FindStudentByID(id int) (*Student, error) {
	// Find student by ID
}

func VerifyCredential(s *Student, cred Credential) bool {
	for _, storedCred := range s.Credentials {
		if storedCred.Hash == cred.Hash {
			return true
		}
	}
}

// function DeleteCredential
func DeleteCredential(s *Student, cred Credential) {
	for i, storedCred := range s.Credentials {
		if storedCred.Hash == cred.Hash {
			s.Credentials = append(s.Credentials[:i], s.Credentials[i+1:]...)
			break
		}
	}
}

// TODO: geenrate Credential Hash
func GenerateCredentialHash(cred Credential) []byte {
	// Generate hash for the credential
}

// TODO: Validata Data Before adding to the blockchain
func ValidateCredentialData(cred Credential) error {
	// Validate credential data
}

// TODO: Credential lookup
// * retrieves and displays all teh credentials of the student
func LookupCredentials(s *Student) {
	// Lookup credentials
}
