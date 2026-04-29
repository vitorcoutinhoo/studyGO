package dto

type CreateColaboradorRequest struct {
	Nome             string  `json:"nome"              binding:"required,min=2,max=100"`
	Email            string  `json:"email"             binding:"required,email,max=100"`
	Telefone         string  `json:"telefone"          binding:"required"`
	Cargo            string  `json:"cargo"             binding:"required"`
	Setor            string  `json:"setor"             binding:"required"`
	Status           string  `json:"status"`
	AtivoPlantao     string  `json:"ativo_plantao"`
	DataAdmissao     string  `json:"data_admissao"     binding:"required"`
	DataDesligamento *string `json:"data_desligamento"`
}

type UpdateColaboradorRequest struct {
	Nome             *string `json:"nome"              binding:"omitempty,min=2,max=100"`
	Email            *string `json:"email"             binding:"omitempty,email,max=100"`
	Telefone         *string `json:"telefone"`
	Cargo            *string `json:"cargo"`
	Setor            *string `json:"setor"`
	Status           *string `json:"status"`
	AtivoPlantao     *string `json:"ativo_plantao"`
	DataAdmissao     *string `json:"data_admissao"`
	DataDesligamento *string `json:"data_desligamento"`
}

type ColaboradorResponse struct {
	Id               string `json:"id"`
	Nome             string `json:"nome"`
	Email            string `json:"email"`
	Telefone         string `json:"telefone"`
	Cargo            string `json:"cargo"`
	Setor            string `json:"setor"`
	Foto             string `json:"foto_url"`
	Status           string `json:"status"`
	AtivoPlantao     string `json:"ativo_plantao"`
	DataAdmissao     string `json:"data_admissao"`
	DataDesligamento string `json:"data_desligamento,omitempty"`
}

type GetColaboradoresByFilterRequest struct {
	Nome         *string `form:"nome"`
	Email        *string `form:"email"`
	Telefone     *string `form:"telefone"`
	Cargo        *string `form:"cargo"`
	Setor        *string `form:"setor"`
	DataAdmissao *string `form:"data_admissao"`
}
