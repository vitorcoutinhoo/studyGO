package dto

type UsuarioRequestDTO struct {
	Email string `json:"email" binding:"required,email,max=100"`
	Senha string `json:"senha" binding:"required,min=6,max=72"`
}

type CadastroByTokenRequestDTO struct {
	Senha string `json:"senha" binding:"required,min=6,max=72"`
}

type UsuarioResponseDTO struct {
	Id            string `json:"id"`
	IdColaborador string `json:"id_colaborador"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Ativo         string `json:"ativo"`
}

type UsuarioAdminRequestDTO struct {
	ColaboradorEmail string `json:"colaborador_email" binding:"required,email,max=100"`
	Email            string `json:"email"             binding:"required,email,max=100"`
	Senha            string `json:"senha"             binding:"required,min=6,max=72"`
	Role             string `json:"role"              binding:"required,oneof=admin gerente colaborador"`
	Ativo            string `json:"ativo"`
}

type LoginRequestDTO struct {
	Email *string `json:"email" binding:"required,email,max=100"`
	Senha *string `json:"senha" binding:"required,min=6,max=72"`
}
