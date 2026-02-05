package dto

type CreateColaboradorRequest struct {
	Nome             string `json:"nome"`
	Email            string `json:"email"`
	Telefone         string `json:"telefone"`
	Cargo            string `json:"cargo"`
	Setor            string `json:"setor"`
	Foto             string `json:"foto"`
	DataAdmissao     string `json:"data_admissao"`
	DataDesligamento string `json:"data_desligamento"`
}

type UpdateColaboradorRequest struct {
	Id               string `json:"id"`
	Email            string `json:"email"`
	Telefone         string `json:"telefone"`
	Cargo            string `json:"cargo"`
	Setor            string `json:"setor"`
	Foto             string `json:"foto"`
	Status           string `json:"status"`
	DataAdmissao     string `json:"data_admissao"`
	DataDesligamento string `json:"data_desligamento"`
}
