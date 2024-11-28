package main

import (
	"log"
	"net/http"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/API/handlers"
	// Add import for blockchain_core
)

func main() {
	// Define routes for student-related operations
	http.HandleFunc("/students", handlers.AddNewStudentAPI)        // Add new student
	http.HandleFunc("/students/update", handlers.UpdateStudentAPI) // Update student details

	// Define a route for Raft-related operations (if needed, e.g., commit transactions)
	http.HandleFunc("/raft/commit", handlers.CommitTransactionAPI) // For Raft log commit

	// Start the API server
	log.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
