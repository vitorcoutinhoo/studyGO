package usuario

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type StatusUsuario int

const (
	StatusAtivo StatusUsuario = iota + 1
	StatusInativo
)

type Role string

const (
	RoleAdmin       = "admin"
	RoleColaborador = "colaborador"
	RoleGerente     = "gerente"
)

type Usuario struct {
	Id            uuid.UUID
	IdColaborador uuid.UUID
	Email         string
	Senha         string
	Role          Role
	Ativo         StatusUsuario
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

var (
	ErrorUserNotFound       = errors.New("Usuário não encontrado!")
	ErrorInvalidEmail       = errors.New("Email Inválido!")
	ErrorEmailexists        = errors.New("Email já existe!")
	ErrorPasswordShort      = errors.New("Senha muito curta, deve conter no minimo 6 caracteres!")
	ErrorInvalidRole        = errors.New("Role Inválida!")
	ErrorInvalidStatus      = errors.New("Status Inválido!")
	ErrorInvalidNome        = errors.New("Nome Inválido!")
	ErrorEmailAlreadyExists = errors.New("Email já existe!")
)

func NewUsuario(idColaborador uuid.UUID, email, senha string, role Role, ativo StatusUsuario) (*Usuario, error) {
	if !isEmailValid(email) {
		return nil, ErrorInvalidEmail
	}

	if len(senha) < 6 {
		return nil, ErrorPasswordShort
	}

	if !isRoleValid(role) {
		return nil, ErrorInvalidRole
	}

	if !isStatusValid(ativo) {
		return nil, ErrorInvalidStatus
	}

	return &Usuario{
		IdColaborador: idColaborador,
		Email:         email,
		Senha:         senha,
		Role:          role,
		Ativo:         ativo,
	}, nil
}

func (u *Usuario) UpdateUsuario(email, senha string, ativo *StatusUsuario) error {
	if email != "" {
		if !isEmailValid(email) {
			return ErrorInvalidEmail
		}

		u.Email = email
	}

	if senha != "" {
		if len(senha) < 6 {
			return ErrorPasswordShort
		}

		u.Senha = senha
	}

	if ativo != nil {
		if !isStatusValid(*ativo) {
			return ErrorInvalidStatus
		}

		u.Ativo = *ativo
	}

	return nil
}

func isEmailValid(email string) bool {
	return len(email) <= 30 && strings.Contains(email, "@")
}

func isRoleValid(role Role) bool {
	switch role {
	case RoleAdmin, RoleColaborador, RoleGerente:
		return true
	}

	return false
}

func isStatusValid(status StatusUsuario) bool {
	return status == StatusAtivo || status == StatusInativo
}
