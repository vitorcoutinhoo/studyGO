package dto

// DTO para criar um colaborador novo
type CreateColaboradorRequest struct {
	Nome             string  `json:"nome"`
	Email            string  `json:"email"`
	Telefone         string  `json:"telefone"`
	Cargo            string  `json:"cargo"`
	Setor            string  `json:"setor"`
	Foto             string  `json:"foto_url"`
	Status           string  `json:"status"`
	AtivoPlantao     string  `json:"ativo_plantao"`
	DataAdmissao     string  `json:"data_admissao"`
	DataDesligamento *string `json:"data_desligamento"`
}

// DTO para atualizar um colaborador
type UpdateColaboradorRequest struct {
	Nome             *string `json:"nome"`
	Email            *string `json:"email"`
	Telefone         *string `json:"telefone"`
	Cargo            *string `json:"cargo"`
	Setor            *string `json:"setor"`
	Foto             *string `json:"foto_url"`
	Status           *string `json:"status"`
	AtivoPlantao     *string `json:"ativo_plantao"`
	DataAdmissao     *string `json:"data_admissao"`
	DataDesligamento *string `json:"data_desligamento"`
}

// Colaborador para retornar dados
type ColaboradorResponse struct {
	Id               string
	Nome             string
	Email            string
	Telefone         string
	Cargo            string
	Setor            string
	Foto             string
	Status           string
	AtivoPlantao     string
	DataAdmissao     string
	DataDesligamento string
}

// DTO para o filtro de pesquisa de colaboradores
// Filtra por Nome, Email, Telefone, Cargo, Setor, DataAdmissao
type GetColaboradoresByFilterRequest struct {
	Nome         *string `form:"nome"`
	Email        *string `form:"email"`
	Telefone     *string `form:"telefone"`
	Cargo        *string `form:"cargo"`
	Setor        *string `form:"setor"`
	DataAdmissao *string `form:"data_admissao"`
}
