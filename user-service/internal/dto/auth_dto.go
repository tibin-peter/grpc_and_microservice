package dto

type RegisterDTO struct {
	ID       uint
	Name     string
	Email    string
	Password string
}

type LoginDTO struct {
	Email    string
	Password string
}

type AuthResponseDTO struct {
	AccessToken  string
	RefreshToken string
}