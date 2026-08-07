package dto

type ModeloComunicacaoRequestDTO struct {
	Nome            string `json:"nome"             binding:"required,min=1,max=255"`
	TipoComunicacao string `json:"tipo_comunicacao" binding:"required"`
	Assunto         string `json:"assunto"          binding:"required,min=1,max=255"`
	Corpo           string `json:"corpo"            binding:"required,min=1"`
}

type ModeloComunicacaoUpdateRequestDTO struct {
	Nome            string `json:"nome"             binding:"omitempty,min=1,max=255"`
	TipoComunicacao string `json:"tipo_comunicacao"`
	Assunto         string `json:"assunto"          binding:"omitempty,min=1,max=255"`
	Corpo           string `json:"corpo"            binding:"omitempty,min=1"`
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
