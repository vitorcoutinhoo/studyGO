package smtp

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"plantao/internal/domain/log"
)

type repositoryFake struct {
	configuracao *Configuracao
	createErr    error
}

func (r *repositoryFake) Get(context.Context) (*Configuracao, error) {
	if r.configuracao == nil {
		return nil, ErrorConfiguracaoSMTPNotFound
	}
	copia := *r.configuracao
	return &copia, nil
}
func (r *repositoryFake) Create(_ context.Context, configuracao *Configuracao) error {
	if r.createErr != nil {
		return r.createErr
	}
	copia := *configuracao
	r.configuracao = &copia
	return nil
}
func (r *repositoryFake) Update(_ context.Context, configuracao *Configuracao) error {
	copia := *configuracao
	r.configuracao = &copia
	return nil
}

type loggerFake struct{}

func (loggerFake) Info(string, ...any)    {}
func (loggerFake) Warn(string, ...any)    {}
func (loggerFake) Error(string, ...any)   {}
func (loggerFake) Debug(string, ...any)   {}
func (loggerFake) With(...any) log.Logger { return loggerFake{} }
func (loggerFake) Fatal(string, ...any)   {}
func (loggerFake) Sync() error            { return nil }

func serviceTeste(repo Repository) *Service { return NewService(repo, loggerFake{}) }

func TestCreateConfiguracaoSMTPValidaCamposECria(t *testing.T) {
	repo := &repositoryFake{}
	configuracao, err := serviceTeste(repo).Create(context.Background(), "smtp.example.com", 587, "usuario", "senha-secreta", "no-reply@example.com")
	if err != nil || configuracao.Id == uuid.Nil || repo.configuracao == nil {
		t.Fatalf("configuração/erro = %+v/%v", configuracao, err)
	}
	for _, teste := range []struct {
		nome, host, usuario, senha, remetente string
		porta                                 int
		esperado                              error
	}{
		{"host", "", "u", "s", "a@example.com", 587, ErrorSMTPHostInvalido},
		{"porta", "smtp.example.com", "u", "s", "a@example.com", 0, ErrorSMTPPortaInvalida},
		{"usuário", "smtp.example.com", "", "s", "a@example.com", 587, ErrorSMTPUsuarioInvalido},
		{"senha", "smtp.example.com", "u", "", "a@example.com", 587, ErrorSMTPPasswordInvalido},
		{"remetente", "smtp.example.com", "u", "s", "inválido", 587, ErrorSMTPFromInvalido},
	} {
		t.Run(teste.nome, func(t *testing.T) {
			_, err := serviceTeste(&repositoryFake{}).Create(context.Background(), teste.host, teste.porta, teste.usuario, teste.senha, teste.remetente)
			if !errors.Is(err, teste.esperado) {
				t.Fatalf("erro = %v", err)
			}
		})
	}
}

func TestUpdateConfiguracaoSMTPPreservaETrocaSenha(t *testing.T) {
	repo := &repositoryFake{configuracao: &Configuracao{Id: uuid.New(), Host: "smtp.example.com", Porta: 587, Usuario: "usuario", Senha: "original", Remetente: "no-reply@example.com"}}
	host := "smtp.novo.example.com"
	configuracao, err := serviceTeste(repo).Update(context.Background(), &AtualizacaoConfiguracao{Host: &host})
	if err != nil || configuracao.Senha != "original" {
		t.Fatalf("senha não preservada: %+v/%v", configuracao, err)
	}
	senha := "nova-senha"
	configuracao, err = serviceTeste(repo).Update(context.Background(), &AtualizacaoConfiguracao{Senha: &senha})
	if err != nil || configuracao.Senha != senha {
		t.Fatalf("senha não atualizada: %+v/%v", configuracao, err)
	}
	_, err = serviceTeste(repo).Update(context.Background(), &AtualizacaoConfiguracao{})
	if !errors.Is(err, ErrorAtualizacaoSMTPVazia) {
		t.Fatalf("erro = %v", err)
	}
}
