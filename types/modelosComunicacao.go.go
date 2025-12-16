package types

import (
	"time"

	"github.com/google/uuid"
)

type ModelosComunicacaoRequest struct {
	Nome      string     `json:"nome"`
	Tipo      string     `json:"tipo"`
	Assunto   string     `json:"assunto"`
	Corpo     string     `json:"corpo"`
	Ativo     bool       `json:"ativo"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ModelosComunicacaoResponse struct {
	ID        string     `json:"id"`
	Nome      string     `json:"nome"`
	Tipo      string     `json:"tipo"`
	Assunto   string     `json:"assunto"`
	Corpo     string     `json:"corpo"`
	Ativo     bool       `json:"ativo"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ModelosComunicacaoInterface interface {
	CreateModelosComunicacao(modelos ModelosComunicacaoRequest) (*ModelosComunicacaoResponse, error)
	GetModelosComunicacao() ([]*ModelosComunicacaoRequest, error)
	GetModelosComunicacaoById(id uuid.UUID) (*ModelosComunicacaoRequest, error)
	UpdateModelosComunicacao(id uuid.UUID, modelos ModelosComunicacaoRequest) error
	DeleteModelosComunicacaoById(id uuid.UUID) error
}
