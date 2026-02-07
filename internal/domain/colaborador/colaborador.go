package colaborador

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type StatusColaborador int

const (
	StatusAtivo StatusColaborador = iota + 1
	StatusInativo
)

// Estrutura do domínio Colaborador
type Colaborador struct {
	Id               uuid.UUID
	Nome             string
	Email            string
	Telefone         string
	Cargo            string
	Setor            string
	Foto             string
	Status           StatusColaborador
	DataAdmissao     *time.Time
	DataDesligamento *time.Time
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

// Erros específicos do domínio Colaborador
var (
	ErrorColaboradorNotFound = errors.New("Colaborador não encontrado!")
	ErrorInvalidEmail        = errors.New("Email Inválido!")
	ErrorinvalidTelefone     = errors.New("Telefone Inválido!")
	ErrorInvalidStatus       = errors.New("Status Inválido!")
)

// Cria um novo colaborador com validações básicas
func NewColaborador(nome, email, telefone, cargo, setor, foto string, dataAdmissao, dataDesligamento *time.Time) (*Colaborador, error) {
	if !isEmailValid(email) {
		return nil, ErrorInvalidEmail
	}

	if !isTelefoneValid(telefone) {
		return nil, ErrorinvalidTelefone
	}

	return &Colaborador{
		Nome:             nome,
		Email:            email,
		Telefone:         telefone,
		Cargo:            cargo,
		Setor:            setor,
		Foto:             foto,
		Status:           StatusAtivo,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
} // Fim NewColaborador

// Atualiza os dados do colaborador com validações básicas
// Só vai atualizar os campos que não forem vazios
func (c *Colaborador) UpdateDados(email, telefone, cargo, setor, foto string, status StatusColaborador) error {
	if email != "" {
		if !isEmailValid(email) {
			return ErrorInvalidEmail
		}

		c.Email = email
	}

	if telefone != "" {
		if !isTelefoneValid(telefone) {
			return ErrorinvalidTelefone
		}

		c.Telefone = telefone
	}

	if status != 0 {
		if !isStatusValid(status) {
			return ErrorInvalidStatus
		}

		c.Status = status
	}

	if cargo != "" {
		c.Cargo = cargo
	}

	if setor != "" {
		c.Setor = setor
	}

	if foto != "" {
		c.Foto = foto
	}

	return nil
} // Fim UpdateDados

// Verifica se o colaborador pode agendar um platão
func (c *Colaborador) PodeAgendarPlatao() (bool, error) {
	if c.Status != StatusAtivo {
		return false, errors.New("Colaborador inativo não pode agendar plantão")
	}

	return true, nil
} // Fim PodeAgendarPlatao

// Valida o email do colaborador
func isEmailValid(email string) bool {
	return len(email) <= 30 && strings.Contains(email, "@")
} // Fim isEmailValid

// Valida o status do colaborador
func isStatusValid(status StatusColaborador) bool {
	return status == StatusAtivo || status == StatusInativo
} // Fim isStatusValid

// Valida o telefone do colaborador
func isTelefoneValid(telefone string) bool {
	var phoneRegexBR = regexp.MustCompile(`^(\+55\s*)?(\(?\d{2}\)?)?\s*9?\s*\d{4}\s*-?\s*\d{4}$`)

	return phoneRegexBR.MatchString(telefone)
} // Fim isTelefoneValid
