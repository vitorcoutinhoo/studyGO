package dto

type ConviteRequestDTO struct {
	IdColaborador    string `json:"id_colaborador"`
	ColaboradorEmail string `json:"colaborador_email"`
	ColaboradorNome  string `json:"colaborador_nome"`
}

type ConviteResponseDTO struct {
	Token         string `json:"token"`
	IdColaborador string `json:"id_colaborador"`
	ExpiraEm      string `json:"expira_em"`
	Usado         bool   `json:"usado"`
	CreatedAt     string `json:"created_at"`
}
