package model

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PatientInfo represents the details of a patient.
type PatientInfo struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	DOB            time.Time `json:"dob"`
	MedicalHistory string    `json:"medical_history"`
}

// Patients holds a collection of patient records.
type Patients struct {
	Records []PatientInfo
}

// NewPatients creates a new Patients instance.
func NewPatients() *Patients {
	return &Patients{
		Records: []PatientInfo{},
	}
}

// AddPatient adds a new patient to the list.
func (p *Patients) AddPatient(patient PatientInfo) {
	p.Records = append(p.Records, patient)
}

// RemovePatient removes a patient from the list by their ID.
func (p *Patients) RemovePatient(id int) error {
	for i, patient := range p.Records {
		if patient.ID == id {
			p.Records = append(p.Records[:i], p.Records[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("patient with ID %d not found", id)
}

// UpdatePatient updates an existing patient's information.
func (p *Patients) UpdatePatient(updatedPatient PatientInfo) error {
	for i, patient := range p.Records {
		if patient.ID == updatedPatient.ID {
			p.Records[i] = updatedPatient
			return nil
		}
	}
	return fmt.Errorf("patient with ID %d not found", updatedPatient.ID)
}

// LoadPatientRecords loads patient records from a JSON file.
func LoadPatientRecords(filePath string) (*Patients, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var records []PatientInfo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&records)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &Patients{Records: records}, nil
}

// SavePatientRecords saves patient records to a JSON file.
func SavePatientRecords(filePath string, patients *Patients) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(patients.Records)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	return nil
}
