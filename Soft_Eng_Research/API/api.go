package main

import (
	"log"
	"net/http"

	"github.com/Research-Group/Soft_Eng_Research/API/handlers"
)

func main() {
	http.HandleFunc("/students", handlers.AddNewStudentAPI) // Example route

	log.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
