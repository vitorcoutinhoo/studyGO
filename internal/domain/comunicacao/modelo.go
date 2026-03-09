package comunicacao

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type StatusModeloComunicacao int

const (
	StatusAtivo StatusModeloComunicacao = iota + 1
	StatusInativo
)

type TipoComunicacao string

const (
	Email TipoComunicacao = "EMAIL"
	SMS   TipoComunicacao = "SMS"
)

type Comunicacao struct {
	Id              uuid.UUID
	Nome            string
	TipoComunicacao TipoComunicacao
	Assunto         string
	Corpo           string
	Ativo           StatusModeloComunicacao
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

var (
	ErrorIvalidNome                     = errors.New("Nome do envio invalido!")
	ErrorInvalidTipoComunicacao         = errors.New("Tipo de comunicação invalido")
	ErrorAssunto                        = errors.New("Assunto invalido")
	ErrorInvalidCorpo                   = errors.New("Corpo da comunicação invalido")
	ErrorInvalidStatus                  = errors.New("Status invalido")
	ErrorModeloComunicacaoAlreadyExists = errors.New("Modelo de comunicação já existe")
	ErrorModeloComunicacaoNotFound      = errors.New("Modelo de comunicação não encontrado")
)

func NewComunicacao(nome, assunto, corpo string, ativo StatusModeloComunicacao, tipoComunicacao TipoComunicacao) (*Comunicacao, error) {
	if len(nome) < 1 {
		return nil, ErrorIvalidNome
	}

	if len(assunto) < 1 {
		return nil, ErrorAssunto
	}

	if len(corpo) < 1 {
		return nil, ErrorInvalidCorpo
	}

	if !isStatusComunicacaoValid(ativo) {
		return nil, ErrorInvalidStatus
	}

	if !isTipoComunicacaoValid(tipoComunicacao) {
		return nil, ErrorInvalidTipoComunicacao
	}

	return &Comunicacao{
		Nome:            nome,
		TipoComunicacao: tipoComunicacao,
		Assunto:         assunto,
		Corpo:           corpo,
		Ativo:           ativo,
	}, nil
}

func (c *Comunicacao) UpdateComunicacao(nome, assunto, corpo *string, ativo *StatusModeloComunicacao, tipoComunicacao *TipoComunicacao) error {
	if nome != nil && *nome != "" {
		c.Nome = *nome
	}

	if tipoComunicacao != nil {
		if !isTipoComunicacaoValid(*tipoComunicacao) {
			return ErrorInvalidTipoComunicacao
		}

		c.TipoComunicacao = *tipoComunicacao
	}

	if assunto != nil && *assunto != "" {
		c.Assunto = *assunto
	}

	if corpo != nil && *corpo != "" {
		c.Corpo = *corpo
	}

	if ativo != nil {
		if !isStatusComunicacaoValid(*ativo) {
			return ErrorInvalidStatus
		}

		c.Ativo = *ativo
	}

	return nil
}

func isStatusComunicacaoValid(status StatusModeloComunicacao) bool {
	return status == StatusAtivo || status == StatusInativo
}

func isTipoComunicacaoValid(t TipoComunicacao) bool {
	switch t {
	case Email, SMS:
		return true
	}

	return false
}
