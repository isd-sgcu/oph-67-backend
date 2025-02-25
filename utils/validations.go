package utils

import (
	"regexp"
)

// IsValidPhone validates a phone number using a regular expression.
func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^((\+66|0)(\d{1,2}\d{3}\d{3,4}))$`)
	return re.MatchString(phone)
}
