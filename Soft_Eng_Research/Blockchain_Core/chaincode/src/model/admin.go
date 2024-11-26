package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Admin struct {
	AdminID string `json:"admin_id"`
	Name    string `json:"name"`
}

func AddNewStudentAPI(w http.ResponseWriter, r *http.Request) {
	// Parse request body to extract student data
	var studentData struct {
		ID        int       `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Age       int       `json:"age"`
		DOB       time.Time `json:"dob"`
		StudentID int       `json:"student_id"`
	}
	json.NewDecoder(r.Body).Decode(&studentData)

	// Create a new instance of the blockchain state
	chain := &StudentChain{}

	// Submit the transaction to the Raft node
	err := node.SubmitTransaction("AddNewStudent", []interface{}{
		studentData.ID,
		studentData.FirstName,
		studentData.LastName,
		studentData.Age,
		studentData.DOB,
		studentData.StudentID,
		chain,
	})
	if err != nil {
		http.Error(w, "Failed to add new student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New student added successfully"))
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
