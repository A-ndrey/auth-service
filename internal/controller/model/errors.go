package model

const EmailError = "email"
const PasswordError = "password"
const CommonError = "common"

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorType string `json:"error_type,omitempty"`
}
