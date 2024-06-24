package helper

import "errors"

type ErrorResponse struct {
	// Message is error message
	Message string `json:"message" validate:"required"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

var errUnknown = errors.New("unknown panic")
