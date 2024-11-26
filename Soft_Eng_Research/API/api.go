package main

import (
	"api/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/students", handlers.AddNewStudentAPI) // Example route

	log.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
