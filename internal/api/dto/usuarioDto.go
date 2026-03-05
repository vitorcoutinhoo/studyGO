package dto

type UsuarioRequestDTO struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required,min=6"`
}

type UsuarioResponseDTO struct {
	Id            string `json:"id"`
	IdColaborador string `json:"id_colaborador"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Ativo         string `json:"ativo"`
}

type UsuarioAdminRequestDTO struct {
	ColaboradorEmail string `json:"colaborador_email" validate:"required,email"`
	Email            string `json:"email" validate:"required,email"`
	Senha            string `json:"senha" validate:"required,min=6"`
	Role             string `json:"role" validate:"required,oneof=admin user"`
	Ativo            string `json:"ativo"`
}
