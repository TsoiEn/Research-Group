package model

import (
	"encoding/json"
	"fmt"
	"os"
)

// PatientRecord represents a patient's record.
type PatientRecord struct {
	ID        int
	Name      string
	Age       int
	Condition string
}

// LoadPatientRecords loads patient records from a JSON file.
func LoadPatientRecords(filePath string) ([]PatientRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var records []PatientRecord
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&records)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return records, nil
}

// SavePatientRecords saves patient records to a JSON file.
func SavePatientRecords(filePath string, records []PatientRecord) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(records)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	return nil
}
