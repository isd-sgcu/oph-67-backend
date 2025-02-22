package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateUID generates a UID in the format AB12345678
func GenerateUID() string {
	// Define letters and digits
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"

	// Seed the random number generator
	rand.NewSource(time.Now().UnixNano())

	// Generate two random letters
	uid := fmt.Sprintf("%c%c", letters[rand.Intn(len(letters))], letters[rand.Intn(len(letters))])

	// Generate eight random digits
	for i := 0; i < 8; i++ {
		uid += string(digits[rand.Intn(len(digits))])
	}

	return uid
}
