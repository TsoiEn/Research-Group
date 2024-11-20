package model

import (
	"strconv"
	"time"
)

type PatientRecord struct {
	ID             int
	Name           string
	DOB            time.Time
	MedicalHistory string
}

type patients struct {
	Patients []PatientRecord
}

func NewPatients() *patients {
	return &patients{
		Patients: []PatientRecord{},
	}
}

func (p *patients) addPatient(patient PatientRecord) {
	p.Patients = append(p.Patients, patient)
}

func (p *patients) findPatient(id string) *PatientRecord {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	for _, patient := range p.Patients {
		if patient.ID == idInt {
			return &patient
		}
	}
	return nil
}
