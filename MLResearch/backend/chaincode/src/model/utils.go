package model

import (
	"fmt"
	"math/big"
	"time"
)

// PatientRecord represents a patient's medical record

// Serialize converts the PatientRecord to a custom byte array format
func (pr *PatientRecord) Serialize() []byte {
	return []byte(fmt.Sprintf("%d|%s|%s|%s", pr.ID, pr.Name, pr.DOB.Format(time.RFC3339), pr.MedicalHistory))
}

// GenerateCredentialHash creates a homomorphic hash of the credential data for integrity
func GenerateCredentialHash(cred *PatientRecord) []byte {
	// Simulated homomorphic hashing using modular arithmetic
	modulus := big.NewInt(1 << 62) // Example modulus for a large space

	// Serialize credential data
	credData := cred.Serialize()
	dataInt := new(big.Int).SetBytes(credData)

	// Homomorphic operation: Modulo addition (simulated homomorphic hashing)
	hash := new(big.Int).Mod(dataInt, modulus)
	return hash.Bytes()
}
