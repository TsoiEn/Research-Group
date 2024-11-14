package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core/chaincode/src/model" // Replace with the actual import path of your `model` package
)

var studentChain = &model.StudentChain{}

// Function to simulate user input for testing admin operations
func testAdminOperations() {

	admin := &model.Admin{
		AdminID: "1",
		Name:    "Admin User",
	}

	// Simulate adding a new student
	fmt.Println("Testing AddNewStudent...")
	newStudent := admin.AddNewStudent(202013432, "John", "Doe", 21, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), 1, studentChain)
	if newStudent != nil {
		studentChain.Students = append(studentChain.Students, newStudent) // Ensure the student is added to the chain
		fmt.Println("AddNewStudent passed:", newStudent)
	} else {
		fmt.Println("AddNewStudent failed.")
	}

	// Simulate adding a credential by admin
	fmt.Println("\nTesting AddCredentialAdmin...")
	cred := model.Credential{
		Type:       model.Academic,
		Issuer:     "Admin University",
		DateIssued: time.Now(),
	}
	adminSuccess := admin.AddCredentialAdmin(newStudent, cred.Type, cred.Issuer, cred.DateIssued)
	if adminSuccess {
		fmt.Println("AddCredentialAdmin passed.")
	} else {
		fmt.Println("AddCredentialAdmin failed.")
	}
}

// Function to simulate user input for testing student operations
func testStudentOperations() {
	student := &model.Student{
		StudentID:   202013432, // Example student ID
		FirstName:   "John",
		LastName:    "Doe",
		Age:         21,
		BirthDate:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Credentials: []*model.Credential{},
	}

	// Simulate adding a credential by a student
	fmt.Println("\nTesting AddCredential...")
	studentID := "202013432"
	cred := model.Credential{
		Type:       model.NonAcademic,
		Issuer:     "Certification Institute",
		DateIssued: time.Now(),
	}
	studentSuccess := student.AddCredential(cred.Type, cred.Issuer, cred.DateIssued)
	if studentSuccess {
		fmt.Println("AddCredential passed.")
	} else {
		fmt.Println("AddCredential failed.")
	}

	// Simulate updating student credentials
	fmt.Println("\nTesting UpdateStudentCredentials...")
	studentIDInt, err := strconv.Atoi(studentID)
	if err != nil {
		fmt.Println("Invalid student ID:", studentID)
		return
	}
	updatedSuccess := studentChain.UpdateStudentCredentials(studentIDInt, cred)
	if updatedSuccess {
		fmt.Println("UpdateStudentCredentials passed.")
	} else {
		fmt.Println("UpdateStudentCredentials failed.")
	}

	// Simulate finding a student by ID
	fmt.Println("\nTesting FindStudentByID...")
	StudentID := studentID
	studentIDInt, err = strconv.Atoi(StudentID)
	if err != nil {
		fmt.Println("Invalid student ID:", StudentID)
		return
	}
	// This should be an existing instance with data
	foundStudent, err := studentChain.FindStudentByID(studentIDInt)
	if err != nil {
		fmt.Println("FindStudentByID failed:", err)
	} else {
		fmt.Printf("FindStudentByID passed. Found student: %v\n", foundStudent)
	}
}

func main() {
	fmt.Println("Running tests for admin and student operations...")

	// Run admin operations tests
	testAdminOperations()

	// Run student operations tests
	testStudentOperations()

	fmt.Println("\nTesting completed.")
}
