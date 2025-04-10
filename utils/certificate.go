package utils

import (
	"crypto/ed25519"
	"encoding/base64"
)

// GenerateED25519Signature generates an ED25519 signature for a given message using the provided private key.
func GenerateED25519Signature(privateKey string, message string) (string, error) {
	priavteKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	if len(priavteKeyBytes) != ed25519.PrivateKeySize {
		return "", err
	}

	signature := ed25519.Sign(priavteKeyBytes, []byte(message))
	return base64.StdEncoding.EncodeToString(signature), nil
}
