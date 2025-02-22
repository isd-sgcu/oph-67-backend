package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// generateTokens creates an access and refresh token
func GenerateTokens(userID string, jwtSecret string) (string, error) {
	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"userId": userID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	access, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return access, nil
}

// DecodeToken decodes the JWT token and returns the userID and any error encountered
func DecodeToken(tokenString string, jwtSecret string) (string, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method to ensure it's using the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract userID from the claims
		userID, ok := claims["userId"].(string)
		if !ok {
			return "", errors.New("userId not found in token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
