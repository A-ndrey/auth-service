package model

const (
	PasswordVeryWeak = iota
	PasswordWeak
	PasswordMedium
	PasswordStrong
	PasswordVeryStrong
)

type PasswordCheckRequest struct {
	Password string `json:"password"`
}

type PasswordCheckResponse struct {
	Strength       int    `json:"strength"`
	MaxStrength    int    `json:"max_strength"`
	MinStrength    int    `json:"min_strength"`
	Recommendation string `json:"recommendation,omitempty"`
}
