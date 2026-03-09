package colaborador

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Constante para o status do colaborador
type StatusColaborador int

const (
	StatusAtivo StatusColaborador = iota + 1
	StatusInativo
)

type CargoColaborador string

const (
	CargoAnalista               CargoColaborador = "Analista"
	CargoGerente                CargoColaborador = "Gerente"
	CargoConsultor              CargoColaborador = "Consultor"
	CargoTecnico                CargoColaborador = "Técnico"
	CargoOutro                  CargoColaborador = "Outro"
	CargoDesenvolvedorFrontend  CargoColaborador = "Desenvolvedor Frontend"
	CargoDesenvolvedorBackend   CargoColaborador = "Desenvolvedor Backend"
	CargoDesenvolvedorFullstack CargoColaborador = "Desenvolvedor Fullstack"
)

type SetorColaborador string

const (
	SetorTI              SetorColaborador = "TI"
	SetorRH              SetorColaborador = "RH"
	SetorFinanceiro      SetorColaborador = "Financeiro"
	SetorSuporte         SetorColaborador = "Suporte"
	SetorDesenvolvimento SetorColaborador = "Desenvolvimento"
	SetorDiretoria       SetorColaborador = "Diretoria"
)

// Estrutura do domínio Colaborador
type Colaborador struct {
	Id               uuid.UUID
	Nome             string
	Email            string
	Telefone         string
	Cargo            CargoColaborador
	Setor            SetorColaborador
	Foto             string
	Status           StatusColaborador
	AtivoPlantao     StatusColaborador
	DataAdmissao     *time.Time
	DataDesligamento *time.Time
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

// Erros específicos do domínio Colaborador
var (
	ErrorColaboradorNotFound = errors.New("Colaborador não encontrado!")
	ErrorInvalidEmail        = errors.New("Email Inválido!")
	ErrorEmailAlreadyExists  = errors.New("Email já existe!")
	ErrorinvalidTelefone     = errors.New("Telefone Inválido!")
	ErrorInvalidStatus       = errors.New("Status Inválido!")
	ErrorInactiveColaborador = errors.New("Colaborador inativo!")
	ErrorInvalidCargo        = errors.New("Cargo inválido!")
	ErrorInvalidSetor        = errors.New("Setor inválido!")
)

// Cria um novo colaborador com validações básicas
func NewColaborador(nome, email, telefone, foto string, dataAdmissao, dataDesligamento *time.Time, ativo, ativoPlatao StatusColaborador, cargo CargoColaborador, setor SetorColaborador) (*Colaborador, error) {
	if !isEmailValid(email) {
		return nil, ErrorInvalidEmail
	}

	if !isTelefoneValid(telefone) {
		return nil, ErrorinvalidTelefone
	}

	if !isStatusValid(ativoPlatao) {
		return nil, ErrorInvalidStatus
	}

	if !isStatusValid(ativo) {
		return nil, ErrorInvalidStatus
	}

	if !isCargoValid(string(cargo)) {
		return nil, ErrorInvalidCargo
	}

	if !isSetorValid(string(setor)) {
		return nil, ErrorInvalidSetor
	}

	return &Colaborador{
		Nome:             nome,
		Email:            email,
		Telefone:         telefone,
		Cargo:            cargo,
		Setor:            setor,
		Foto:             foto,
		Status:           StatusAtivo,
		AtivoPlantao:     ativoPlatao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
} // Fim NewColaborador

// Atualiza os dados do colaborador com validações básicas
// Só vai atualizar os campos que não forem vazios
func (c *Colaborador) UpdateDados(nome, email, telefone, foto *string, dataAdmissao, dataDesligamento *time.Time, status, ativoPlantao *StatusColaborador, cargo *CargoColaborador, setor *SetorColaborador) error {
	if nome != nil && *nome != "" {
		c.Nome = *nome
	}

	if email != nil && *email != "" {
		if !isEmailValid(*email) {
			return ErrorInvalidEmail
		}

		c.Email = *email
	}

	if telefone != nil && *telefone != "" {
		if !isTelefoneValid(*telefone) {
			return ErrorinvalidTelefone
		}

		c.Telefone = *telefone
	}

	if status != nil && *status != 0 {
		if !isStatusValid(*status) {
			return ErrorInvalidStatus
		}

		c.Status = *status
	}

	if ativoPlantao != nil && *ativoPlantao != 0 {
		if !isStatusValid(*ativoPlantao) {
			return ErrorInvalidStatus
		}

		c.AtivoPlantao = *ativoPlantao
	}

	if cargo != nil && *cargo != "" {
		c.Cargo = *cargo
	}

	if setor != nil && *setor != "" {
		c.Setor = *setor
	}

	if foto != nil && *foto != "" {
		c.Foto = *foto
	}

	if dataAdmissao != nil {
		c.DataAdmissao = dataAdmissao
	}

	if dataDesligamento != nil {
		c.DataDesligamento = dataDesligamento
	}

	return nil
} // Fim UpdateDados

// Verifica se o colaborador pode agendar um platão
func (c *Colaborador) PodeAgendarPlatao() (bool, error) {
	if c.AtivoPlantao != StatusAtivo {
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

func isCargoValid(cargo string) bool {
	switch cargo {
	case string(CargoAnalista), string(CargoGerente), string(CargoConsultor), string(CargoTecnico), string(CargoOutro), string(CargoDesenvolvedorFrontend), string(CargoDesenvolvedorBackend), string(CargoDesenvolvedorFullstack):
		return true
	default:
		return false
	}
}

func isSetorValid(setor string) bool {
	switch setor {
	case string(SetorTI), string(SetorRH), string(SetorFinanceiro), string(SetorSuporte), string(SetorDesenvolvimento), string(SetorDiretoria):
		return true
	default:
		return false
	}
}
