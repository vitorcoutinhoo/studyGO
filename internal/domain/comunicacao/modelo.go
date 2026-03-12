package comunicacao

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type StatusModeloComunicacao int

const (
	StatusAtivo StatusModeloComunicacao = iota + 1
	StatusInativo
)

type TipoComunicacao string

const (
	PlantaoAgendado       TipoComunicacao = "Plantão Agendado"
	PlantaoConluido       TipoComunicacao = "Plantão Concluido"
	PlantaoAindaAberto    TipoComunicacao = "Plantão Ainda Está Aberto"
	PlantaoPago           TipoComunicacao = "Plantão Pago"
	ColaboradorCadastrado TipoComunicacao = "Colaborador Cadastrado"
	ColaboradorAtualizado TipoComunicacao = "Colaborador Atualizado"
	ColaboradorDeletado   TipoComunicacao = "Colaborador Deletado"
	UsuarioCadastrado     TipoComunicacao = "Usuário Cadastrado"
	EmailAtualizado       TipoComunicacao = "Email do Usuário Atualizado"
	SenhaAtualizada       TipoComunicacao = "Senha do Usuário Atualizada"
	UsuarioDeletado       TipoComunicacao = "Usuário Deletado"
)

type TagBody string

const (
	Nome       TagBody = "nome"
	DataInicio TagBody = "dataInicio"
	DataFim    TagBody = "dataFim"
	ValorPago  TagBody = "valorPago"
	Email      TagBody = "email"
	DataAtual  TagBody = "dataAtual"
)

var requiredTags = map[TipoComunicacao][]TagBody{
	PlantaoAgendado:       {Nome, DataInicio, DataFim},
	PlantaoConluido:       {Nome, DataInicio, DataFim},
	PlantaoAindaAberto:    {Nome, DataInicio, DataFim},
	PlantaoPago:           {Nome, DataInicio, DataFim, ValorPago},
	ColaboradorCadastrado: {Nome, Email},
	ColaboradorAtualizado: {Nome, Email},
	ColaboradorDeletado:   {Nome, DataAtual},
	UsuarioCadastrado:     {Nome, Email},
	EmailAtualizado:       {Nome, Email},
	SenhaAtualizada:       {Nome, DataAtual},
	UsuarioDeletado:       {Nome, DataAtual},
}

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
	ErrorInvalidNome                    = errors.New("Nome do envio invalido!")
	ErrorInvalidTipoComunicacao         = errors.New("Tipo de comunicação invalido")
	ErrorAssunto                        = errors.New("Assunto invalido")
	ErrorInvalidCorpo                   = errors.New("Corpo da comunicação invalido")
	ErrorInvalidStatus                  = errors.New("Status invalido")
	ErrorModeloComunicacaoAlreadyExists = errors.New("Modelo de comunicação já existe")
	ErrorModeloComunicacaoNotFound      = errors.New("Modelo de comunicação não encontrado")
)

func NewComunicacao(nome, assunto, corpo string, ativo StatusModeloComunicacao, tipoComunicacao TipoComunicacao) (*Comunicacao, error) {
	if len(nome) < 1 {
		return nil, ErrorInvalidNome
	}

	if len(assunto) < 1 {
		return nil, ErrorAssunto
	}

	err := validateEmailBodyTag(tipoComunicacao, corpo)

	if err != nil {
		return nil, err
	}

	err = isValidHTML(corpo)

	if err != nil {
		return nil, err
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

	if corpo != nil {
		err := validateEmailBodyTag(*tipoComunicacao, *corpo)

		if err != nil {
			return err
		}

		err = isValidHTML(*corpo)

		if err != nil {
			return err
		}

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
	case
		PlantaoAgendado,
		PlantaoConluido,
		PlantaoAindaAberto,
		PlantaoPago,
		ColaboradorCadastrado,
		ColaboradorAtualizado,
		ColaboradorDeletado,
		UsuarioCadastrado,
		EmailAtualizado,
		SenhaAtualizada,
		UsuarioDeletado:
		return true
	}

	return false
}

func validateEmailBodyTag(tipoComunicacao TipoComunicacao, body string) error {
	tags, ok := requiredTags[tipoComunicacao]

	if !ok {
		return fmt.Errorf("tipo de comunicação inválido: %s", tipoComunicacao)
	}

	var missingTags []string

	for _, tag := range tags {
		pattern := fmt.Sprintf(`{{\s*\.%s\s*}}`, tag)
		matched, _ := regexp.MatchString(pattern, body)

		if !matched {
			missingTags = append(missingTags, fmt.Sprintf("{{%s}}", tag))
		}
	}

	if len(missingTags) > 0 {
		return fmt.Errorf(
			"tags obrigatórias ausentes para o tipo '%s': %s",
			tipoComunicacao,
			strings.Join(missingTags, ", "),
		)
	}

	return nil
}

func isValidHTML(htmlBody string) error {
	if strings.TrimSpace(htmlBody) == "" {
		return errors.New("corpo do email não pode estar vazio")
	}

	_, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return errors.New("HTML inválido no corpo do email")
	}

	if strings.Contains(strings.ToLower(htmlBody), "<script") {
		return errors.New("scripts não são permitidos no email")
	}

	return nil
}
