package dto

type AuthenticateRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
