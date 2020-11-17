package model

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DefineUserResponse struct {
	Email          string `json:"email"`
	TokenExpiresAt int64  `json:"token_expires_at"`
}
