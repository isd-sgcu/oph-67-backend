package domain

import "errors"

// ErrorResponse represents a basic error structure
type ErrorResponse struct {
	Error   string  `json:"error"`
	Message *string `json:"message,omitempty"`
}

var ErrUserAlreadyEntered = errors.New("user has already entered")
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyStaff = errors.New("user is already a staff")
var ErrUserNotCentralStaff = errors.New("user is not a central staff")
var ErrStudentEvaluationAlreadyExists = errors.New("evaluation for this student already exists")
