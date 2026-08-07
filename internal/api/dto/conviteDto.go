package dto

type ConviteRequestDTO struct {
	IdColaborador    string `json:"id_colaborador"    binding:"required,uuid"`
	ColaboradorEmail string `json:"colaborador_email" binding:"required,email,max=100"`
	ColaboradorNome  string `json:"colaborador_nome"  binding:"required,min=2,max=100"`
}

type ConviteResponseDTO struct {
	Token         string `json:"token"`
	IdColaborador string `json:"id_colaborador"`
	ExpiraEm      string `json:"expira_em"`
	Usado         bool   `json:"usado"`
	CreatedAt     string `json:"created_at"`
}
