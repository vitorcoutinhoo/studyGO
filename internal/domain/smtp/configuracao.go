package smtp

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"plantao/internal/domain/log"
)

var (
	ErrorConfiguracaoSMTPNotFound      = errors.New("configuração SMTP não encontrada")
	ErrorConfiguracaoSMTPAlreadyExists = errors.New("configuração SMTP já existe")
	ErrorSMTPHostInvalido              = errors.New("host SMTP é obrigatório")
	ErrorSMTPPortaInvalida             = errors.New("porta SMTP deve estar entre 1 e 65535")
	ErrorSMTPUsuarioInvalido           = errors.New("usuário SMTP é obrigatório")
	ErrorSMTPPasswordInvalido          = errors.New("senha SMTP é obrigatória")
	ErrorSMTPFromInvalido              = errors.New("remetente SMTP inválido")
	ErrorAtualizacaoSMTPVazia          = errors.New("informe ao menos um campo para atualização")
)

type Configuracao struct {
	Id        uuid.UUID
	Host      string
	Porta     int
	Usuario   string
	Senha     string
	Remetente string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AtualizacaoConfiguracao struct {
	Host      *string
	Porta     *int
	Usuario   *string
	Senha     *string
	Remetente *string
}

type Repository interface {
	Get(ctx context.Context) (*Configuracao, error)
	Create(ctx context.Context, configuracao *Configuracao) error
	Update(ctx context.Context, configuracao *Configuracao) error
}

type Service struct {
	repository Repository
	log        log.Logger
}

func NewService(repository Repository, log log.Logger) *Service {
	return &Service{repository: repository, log: log}
}

func (s *Service) Get(ctx context.Context) (*Configuracao, error) {
	configuracao, err := s.repository.Get(ctx)
	if err != nil {
		return nil, err
	}
	return configuracao, nil
}

func (s *Service) Create(ctx context.Context, host string, porta int, usuario, senha, remetente string) (*Configuracao, error) {
	configuracao := &Configuracao{Id: uuid.New(), Host: host, Porta: porta, Usuario: usuario, Senha: senha, Remetente: remetente}
	if err := validar(configuracao); err != nil {
		return nil, err
	}
	if err := s.repository.Create(ctx, configuracao); err != nil {
		return nil, err
	}
	s.log.Info("configuração SMTP criada", "id_configuracao_smtp", configuracao.Id)
	return configuracao, nil
}

func (s *Service) Update(ctx context.Context, atualizacao *AtualizacaoConfiguracao) (*Configuracao, error) {
	if atualizacao == nil || !atualizacao.temCampos() {
		return nil, ErrorAtualizacaoSMTPVazia
	}
	configuracao, err := s.repository.Get(ctx)
	if err != nil {
		return nil, err
	}
	if atualizacao.Host != nil {
		configuracao.Host = *atualizacao.Host
	}
	if atualizacao.Porta != nil {
		configuracao.Porta = *atualizacao.Porta
	}
	if atualizacao.Usuario != nil {
		configuracao.Usuario = *atualizacao.Usuario
	}
	if atualizacao.Senha != nil {
		configuracao.Senha = *atualizacao.Senha
	}
	if atualizacao.Remetente != nil {
		configuracao.Remetente = *atualizacao.Remetente
	}
	if err := validar(configuracao); err != nil {
		return nil, err
	}
	if err := s.repository.Update(ctx, configuracao); err != nil {
		return nil, err
	}
	s.log.Info("configuração SMTP atualizada", "id_configuracao_smtp", configuracao.Id)
	return configuracao, nil
}

func (a *AtualizacaoConfiguracao) temCampos() bool {
	return a.Host != nil || a.Porta != nil || a.Usuario != nil || a.Senha != nil || a.Remetente != nil
}

func validar(configuracao *Configuracao) error {
	if strings.TrimSpace(configuracao.Host) == "" {
		return ErrorSMTPHostInvalido
	}
	if configuracao.Porta < 1 || configuracao.Porta > 65535 {
		return ErrorSMTPPortaInvalida
	}
	if strings.TrimSpace(configuracao.Usuario) == "" {
		return ErrorSMTPUsuarioInvalido
	}
	if strings.TrimSpace(configuracao.Senha) == "" {
		return ErrorSMTPPasswordInvalido
	}
	remetente := strings.TrimSpace(configuracao.Remetente)
	endereco, err := mail.ParseAddress(remetente)
	if err != nil || endereco.Address != remetente {
		return ErrorSMTPFromInvalido
	}
	return nil
}
