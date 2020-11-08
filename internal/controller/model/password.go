package model

const (
	PasswordVeryWeak   = "very weak"
	PasswordWeak       = "weak"
	PasswordMedium     = "medium"
	PasswordStrong     = "strong"
	PasswordVeryStrong = "very strong"
)

type PasswordCheckRequest struct {
	Password string `json:"password"`
}

type PasswordCheckResponse struct {
	Strength       string `json:"strength"`
	Recommendation string `json:"recommendation,omitempty"`
}
