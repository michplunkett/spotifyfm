package utility

import (
	"crypto/rand"
	"encoding/base64"
)

type HelperFunctions interface {
	ArrayHasNoEmptyStrings(envVars []string) bool
	generateRandomBytes(n int) ([]byte, error)
	GenerateRandomString(s int) (string, error)
}

func ArrayHasNoEmptyStrings(envVars []string) bool {
	for _, value := range envVars {
		if value == EmptyString {
			return false
		}
	}

	return true
}

func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
