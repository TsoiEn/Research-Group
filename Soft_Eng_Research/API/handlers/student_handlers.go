package handlers

import (
	"encoding/json"
	"net/http"
	"time"
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

	// Interact with the blockchain via node.SubmitTransaction
	err = node.SubmitTransaction("AddNewStudent", []interface{}{
		studentData.ID,
		studentData.FirstName,
		studentData.LastName,
		studentData.Age,
		studentData.DOB,
		studentData.StudentID,
	})
	if err != nil {
		http.Error(w, "Failed to add new student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New student added successfully"))
}
