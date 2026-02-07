package dto

type CreateColaboradorRequest struct {
	Nome             string `json:"nome"`
	Email            string `json:"email"`
	Telefone         string `json:"telefone"`
	Cargo            string `json:"cargo"`
	Setor            string `json:"setor"`
	Foto             string `json:"foto_url"`
	DataAdmissao     string `json:"data_admissao"`
	DataDesligamento string `json:"data_desligamento"`
}

type UpdateColaboradorRequest struct {
	Email            string  `json:"email"`
	Telefone         string  `json:"telefone"`
	Cargo            string  `json:"cargo"`
	Setor            string  `json:"setor"`
	Foto             string  `json:"foto"`
	Status           string  `json:"status"`
	DataAdmissao     string  `json:"data_admissao"`
	DataDesligamento *string `json:"data_desligamento"`
}

type GetColaboradoresByFilterRequest struct {
	Nome         *string `form:"nome"`
	Email        *string `form:"email"`
	Telefone     *string `form:"telefone"`
	Cargo        *string `form:"cargo"`
	Setor        *string `form:"setor"`
	DataAdmissao *string `form:"data_admissao"`
}
