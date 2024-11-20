package model

import (
	"fmt"
	"math/big"
	"time"
)

// Serialize converts the PatientRecord to a custom byte array format
func (pr *PatientInfo) Serialize() []byte {
	return []byte(fmt.Sprintf("%d|%s|%s|%s", pr.ID, pr.Name, pr.DOB.Format(time.RFC3339), pr.MedicalHistory))
}

// GenerateCredentialHash creates a homomorphic hash of the credential data for integrity
func GenerateCredentialHash(cred *PatientInfo) []byte {

	modulus := big.NewInt(1 << 62)

	credData := cred.Serialize()
	dataInt := new(big.Int).SetBytes(credData)

	hash := new(big.Int).Mod(dataInt, modulus)
	return hash.Bytes()
}
