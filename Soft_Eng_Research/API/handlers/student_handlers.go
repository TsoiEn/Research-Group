package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/consensus"
)

func AddNewStudentAPI(w http.ResponseWriter, r *http.Request) {
	var studentData struct {
		ID        int       `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Age       int       `json:"age"`
		DOB       time.Time `json:"dob"`
		StudentID int       `json:"student_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&studentData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Create the transaction payload
	transaction := map[string]interface{}{
		"action": "AddNewStudent",
		"data": []interface{}{
			studentData.ID,
			studentData.FirstName,
			studentData.LastName,
			studentData.Age,
			studentData.DOB,
			studentData.StudentID,
		},
	}

	// Create a new Raft Node (you can expand this to handle more nodes and state)
	node := consensus.CreateRaftNode("node1", 5*time.Second, 10*time.Second)
	fmt.Println(node)

	// Submit the transaction
	err = node.ProposeTransaction(transaction)
	if err != nil {
		http.Error(w, "Failed to add new student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New student added successfully"))
}
func UpdateStudentAPI(w http.ResponseWriter, r *http.Request) {
	var studentData struct {
		ID        int       `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Age       int       `json:"age"`
		DOB       time.Time `json:"dob"`
		StudentID int       `json:"student_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&studentData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Create the transaction payload
	transaction := map[string]interface{}{
		"action": "UpdateStudent",
		"data": []interface{}{
			studentData.ID,
			studentData.FirstName,
			studentData.LastName,
			studentData.Age,
			studentData.DOB,
			studentData.StudentID,
		},
	}

	// Create a new Raft Node (you can expand this to handle more nodes and state)
	node := consensus.CreateRaftNode("node1", 5*time.Second, 10*time.Second)
	fmt.Println(node)

	// Submit the transaction
	err = node.ProposeTransaction(transaction)
	if err != nil {
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Student updated successfully"))
}

func CommitTransactionAPI(w http.ResponseWriter, r *http.Request) {
	var transactionData struct {
		TransactionID string `json:"transaction_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&transactionData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Create a new Raft Node (you can expand this to handle more nodes and state)
	node := consensus.CreateRaftNode("node1", 5*time.Second, 10*time.Second)
	fmt.Println(node)

	// Commit the transaction
	err = node.CommitTransaction(transactionData.TransactionID)
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction committed successfully"))
}
