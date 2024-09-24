package chaincode

import (
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

// CRUD operations for Student

func (c *Chaincode) AddStudent(ctx contractapi.TransactionContextInterface, student Student) error {
	exist, err := c.StudentExists(ctx, student.ID)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("the student %s already exists", student.ID)
	}

	Stnd := &Student{
		ID:        student.ID,
		FirstName: student.FirstName,
		LastName:  student.LastName,
		Age:       student.Age,
		BirthDate: student.BirthDate,
	}
	studentJSON, err := json.Marshal(Stnd)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(student.ID, studentJSON)

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

func (c *Chaincode) UpdateStudent(ctx contractapi.TransactionContextInterface, student Student) error {
	exist, err := c.StudentExists(ctx, student.ID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("the student %s does not exist", student.ID)
	}

	Stnd := &Student{
		ID:        student.ID,
		FirstName: student.FirstName,
		LastName:  student.LastName,
		Age:       student.Age,
		BirthDate: student.BirthDate,
	}
	studentJSON, err := json.Marshal(Stnd)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(student.ID, studentJSON)
}

func (c *Chaincode) DeleteStudent(ctx contractapi.TransactionContextInterface, student Student) error {
	exist, err := c.StudentExists(ctx, student.ID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("the student %s does not exist", student.ID)
	}

	stnd := &Student{
		ID: student.ID,
	}

	return ctx.GetStub().DelState(stnd.ID)
}

func (c *Chaincode) StudentExists(ctx contractapi.TransactionContextInterface, studentID string) (bool, error) {
	studentJSON, err := ctx.GetStub().GetState(studentID)
	if err != nil {
		return false, err
	}

	return studentJSON != nil, nil
}
