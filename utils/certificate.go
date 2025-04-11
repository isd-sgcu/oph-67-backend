package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
)

// CreateToken creates a token(URLsafe) using the provided message and private key
func GenerateED25519Signature(privateKey string, message string) (string, error) {
	priavteKeyBytes, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	if len(priavteKeyBytes) != ed25519.PrivateKeySize {
		return "", errors.New("invalid private key size")
	}

	signature := ed25519.Sign(priavteKeyBytes, []byte(message))
	encodedMessage := base64.RawURLEncoding.EncodeToString([]byte(message))

	// Combine the encoded message and signature
	token := encodedMessage + "." + base64.RawURLEncoding.EncodeToString(signature)

	return token, nil
}
