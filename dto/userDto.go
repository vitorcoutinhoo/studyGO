package dto

type UserCreateDTO struct {
	Email     string
	SenhaHash string
}

type UserResponseDTO struct {
	ID    string
	Email string
	Role  string
}
