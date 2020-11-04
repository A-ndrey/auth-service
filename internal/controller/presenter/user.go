package presenter

type UserRequest struct {
	Service  string `json:"service"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
