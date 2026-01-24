package colaborador

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status map[string]int

var (
	ColaboradorAtivo   = Status{"Ativo": 1}
	ColaboradorInativo = Status{"Inativo": 2}
)

type Colaborador struct {
	Id        *uuid.UUID
	Nome      string
	Email     string
	Telefone  string
	Cargo     string
	Setor     string
	Foto      string
	Status    Status
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

var (
	ErrorInvalidEmail    = "Email Inválido!"
	ErrorinvalidTelefone = "Telefone Inválido!"
	ErrorInvalidStatus   = "Status Inválido!"
)

func NewColaborador(nome, email, telefone, cargo, setor, foto string) (*Colaborador, error) {

	return &Colaborador{
		Nome:     nome,
		Email:    email,
		Telefone: telefone,
		Cargo:    cargo,
		Setor:    setor,
		Foto:     foto,
		Status:   ColaboradorAtivo,
	}, nil
}

func isEmailValid(email string) bool {
	return len(email) <= 30 && strings.Contains(email, "@")
}

func isStatusValid(status Status) bool {
	return false
}
