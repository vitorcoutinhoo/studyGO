package dto

type ModeloComunicacaoRequestDTO struct {
	Nome            string `json:"nome"`
	TipoComunicacao string `json:"tipo_comunicacao"`
	Assunto         string `json:"assunto"`
	Corpo           string `json:"corpo"`
}

type ModeloComunicacaoUpdateRequestDTO struct {
	Nome            string `json:"nome"`
	TipoComunicacao string `json:"tipo_comunicacao"`
	Assunto         string `json:"assunto"`
	Corpo           string `json:"corpo"`
	Ativo           string `json:"ativo"`
}

type ModeloComunicacaoResponseDTO struct {
	Id              string `json:"id"`
	Nome            string `json:"nome"`
	TipoComunicacao string `json:"tipo_comunicacao"`
	Assunto         string `json:"assunto"`
	Corpo           string `json:"corpo"`
	Ativo           string `json:"ativo"`
}
