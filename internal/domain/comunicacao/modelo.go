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

type NomeEnvio string

const (
	BoasVindas         NomeEnvio = "Boas Vindas"
	NovoPlantao        NomeEnvio = "Novo Plantão"
	CadastroAtualizado NomeEnvio = "Cadastro Atualizado"
	CadastroExcluido   NomeEnvio = "Cadastro Excluido"
)

type Comunicacao struct {
	Id              uuid.UUID
	Nome            NomeEnvio
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

func NewComunicacao(assunto, corpo string, ativo StatusModeloComunicacao, tipoComunicacao TipoComunicacao, nome NomeEnvio) (*Comunicacao, error) {
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

	if !isNomeEnvioValid(nome) {
		return nil, ErrorIvalidNome
	}

	return &Comunicacao{
		Nome:            nome,
		TipoComunicacao: tipoComunicacao,
		Assunto:         assunto,
		Corpo:           corpo,
		Ativo:           ativo,
	}, nil
}

func (c *Comunicacao) UpdateComunicacao(assunto, corpo string, ativo *StatusModeloComunicacao, tipoComunicacao *TipoComunicacao, nome *NomeEnvio) error {
	if nome != nil {
		if !isNomeEnvioValid(*nome) {
			return ErrorIvalidNome
		}
		c.Nome = *nome
	}

	if tipoComunicacao != nil {
		if !isTipoComunicacaoValid(*tipoComunicacao) {
			return ErrorInvalidTipoComunicacao
		}

		c.TipoComunicacao = *tipoComunicacao
	}

	if assunto != "" {
		c.Assunto = assunto
	}

	if corpo != "" {
		c.Corpo = corpo
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

func isNomeEnvioValid(n NomeEnvio) bool {
	switch n {
	case BoasVindas, NovoPlantao, CadastroAtualizado, CadastroExcluido:
		return true
	}

	return false
}
