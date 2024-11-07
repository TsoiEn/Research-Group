package chaincode

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Chaincode definition
type Student struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	BirthDate string `json:"birthDate"`
}

type Chaincode struct {
	contractapi.Contract
}

// Credential Handling Functions

// Verifies if a given credential for a student is valid by checking if it matches any stored credential hashes.
func (c *Chaincode) VerifyCredential(ctx contractapi.TransactionContextInterface, studentID string, credentialData string) (bool, error) {
	// Retrieve the student from the blockchain
	_, err := c.ReadStudent(ctx, studentID)
	if err != nil {
		return false, err
	}

	// Hash the provided credential data
	hashedCredential := hashCredential(credentialData)

	// Retrieve stored credentials for the student
	storedCredentials, err := ctx.GetStub().GetState(studentID + "_credentials")
	if err != nil {
		return false, err
	}
	if storedCredentials == nil {
		return false, fmt.Errorf("no credentials found for student %s", studentID)
	}

	// Unmarshal stored credentials
	var credentials []string
	if err = json.Unmarshal(storedCredentials, &credentials); err != nil {
		return false, err
	}

	// Compare the hash of the provided credential with stored credentials
	for _, storedHash := range credentials {
		if hashedCredential == storedHash {
			return true, nil
		}
	}

	return false, fmt.Errorf("credential is not valid")
}

// Hashes the credential data
func hashCredential(credentialData string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(credentialData)))
}

// Allows the admin to update a specific credential for a student, either by modifying the credential’s properties or adding new credentials, depending on project needs.
func (c *Chaincode) UpdateCredential(ctx contractapi.TransactionContextInterface, studentID string, credentialData string) error {
	// Check if the student exists
	exists, err := c.StudentExists(ctx, studentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the student %s does not exist", studentID)
	}

	// Hash the new credential data
	hashedCredential := hashCredential(credentialData)

	// Retrieve stored credentials for the student
	storedCredentials, err := ctx.GetStub().GetState(studentID + "_credentials")
	if err != nil {
		return err
	}

	// Unmarshal stored credentials
	var credentials []string
	if storedCredentials != nil {
		if err = json.Unmarshal(storedCredentials, &credentials); err != nil {
			return err
		}
	}

	// Add the new hashed credential to the list
	credentials = append(credentials, hashedCredential)

	// Marshal the updated credentials
	updatedCredentials, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	// Save the updated credentials back to the blockchain
	return ctx.GetStub().PutState(studentID+"_credentials", updatedCredentials)
}

// Retrieves all credentials associated with a specific student ID from the blockchain.
func (c *Chaincode) RetrieveCredential(ctx contractapi.TransactionContextInterface, studentID string) ([]string, error) {
	// Check if the student exists
	exists, err := c.StudentExists(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("the student %s does not exist", studentID)
	}

	// Retrieve stored credentials for the student
	storedCredentials, err := ctx.GetStub().GetState(studentID + "_credentials")
	if err != nil {
		return nil, err
	}
	if storedCredentials == nil {
		return nil, fmt.Errorf("no credentials found for student %s", studentID)
	}

	// Unmarshal stored credentials
	var credentials []string
	err = json.Unmarshal(storedCredentials, &credentials)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (c *Chaincode) ReadStudent(ctx contractapi.TransactionContextInterface, studentID string) (*Student, error) {
	studentJSON, err := ctx.GetStub().GetState(studentID)
	if err != nil {
		return nil, err
	}
	if studentJSON == nil {
		return nil, fmt.Errorf("the student %s does not exist", studentID)
	}

	var student Student
	err = json.Unmarshal(studentJSON, &student)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (c *Chaincode) StudentExists(ctx contractapi.TransactionContextInterface, studentID string) (bool, error) {
	studentJSON, err := ctx.GetStub().GetState(studentID)
	if err != nil {
		return false, err
	}

	return studentJSON != nil, nil
}

// UpdateStudent allows the admin to update a student's details
func (c *Chaincode) UpdateStudent(ctx contractapi.TransactionContextInterface, student Student) error {
	exist, err := c.StudentExists(ctx, student.ID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("the student %s does not exist", student.ID)
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(student.ID, studentJSON)
}

// AddAcademicCredential allows the admin to add academic credentials for a specific student.
func (c *Chaincode) AddAcademicCredential(ctx contractapi.TransactionContextInterface, studentID string, credentialData string) error {
	// Check if the student exists
	exists, err := c.StudentExists(ctx, studentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the student %s does not exist", studentID)
	}

	// Hash the academic credential data
	hashedCredential := hashCredential(credentialData)

	// Retrieve stored academic credentials for the student
	storedCredentials, err := ctx.GetStub().GetState(studentID + "_academic_credentials")
	if err != nil {
		return err
	}

	// Unmarshal stored credentials
	var credentials []string
	if storedCredentials != nil {
		err = json.Unmarshal(storedCredentials, &credentials)
		if err != nil {
			return err
		}
	}

	// Check if the credential already exists, to avoid duplicates
	for _, existing := range credentials {
		if existing == hashedCredential {
			return fmt.Errorf("the credential already exists for student %s", studentID)
		}
	}

	// Add the new hashed academic credential to the list
	credentials = append(credentials, hashedCredential)

	// Marshal the updated credentials
	updatedCredentials, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	// Save the updated credentials back to the blockchain
	return ctx.GetStub().PutState(studentID+"_academic_credentials", updatedCredentials)
}

// UpdateAcademicCredential allows the admin to update existing academic credentials for a specific student.
func (c *Chaincode) UpdateAcademicCredential(ctx contractapi.TransactionContextInterface, studentID string, oldCredentialData string, newCredentialData string) error {
	// Check if the student exists
	exists, err := c.StudentExists(ctx, studentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the student %s does not exist", studentID)
	}

	// Hash the old and new credential data
	hashedOldCredential := hashCredential(oldCredentialData)
	hashedNewCredential := hashCredential(newCredentialData)

	// Retrieve stored academic credentials for the student
	storedCredentials, err := ctx.GetStub().GetState(studentID + "_academic_credentials")
	if err != nil {
		return err
	}

	// Unmarshal stored credentials
	var credentials []string
	if storedCredentials != nil {
		err = json.Unmarshal(storedCredentials, &credentials)
		if err != nil {
			return err
		}
	}

	// Check if the old credential exists
	updated := false
	for i, storedHash := range credentials {
		if storedHash == hashedOldCredential {
			// Replace with new credential
			credentials[i] = hashedNewCredential
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("old credential not found for student %s", studentID)
	}

	// Marshal the updated credentials
	updatedCredentials, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	// Save the updated credentials back to the blockchain
	return ctx.GetStub().PutState(studentID+"_academic_credentials", updatedCredentials)
}

// AddNonAcademicCredential allows students to add non-academic credentials for themselves.
func (c *Chaincode) AddNonAcademicCredential(ctx contractapi.TransactionContextInterface, studentID string, credentialData string) error {
	// Check if the student exists
	exists, err := c.StudentExists(ctx, studentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the student %s does not exist", studentID)
	}

	// Hash the non-academic credential data
	hashedCredential := hashCredential(credentialData)

	// Retrieve stored non-academic credentials for the student
	storedCredentials, err := ctx.GetStub().GetState(studentID + "_non_academic_credentials")
	if err != nil {
		return err
	}

	// Unmarshal stored credentials
	var credentials []string
	if storedCredentials != nil {
		err = json.Unmarshal(storedCredentials, &credentials)
		if err != nil {
			return err
		}
	}

	// Add the new hashed non-academic credential to the list
	credentials = append(credentials, hashedCredential)

	// Marshal the updated credentials
	updatedCredentials, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	// Save the updated credentials back to the blockchain
	return ctx.GetStub().PutState(studentID+"_non_academic_credentials", updatedCredentials)
}
