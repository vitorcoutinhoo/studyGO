package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"plantao/internal/domain/log"
	"plantao/internal/domain/smtp"
)

type smtpRepositoryControllerFake struct{ configuracao *smtp.Configuracao }

func (r *smtpRepositoryControllerFake) Get(context.Context) (*smtp.Configuracao, error) {
	if r.configuracao == nil {
		return nil, smtp.ErrorConfiguracaoSMTPNotFound
	}
	copia := *r.configuracao
	return &copia, nil
}
func (r *smtpRepositoryControllerFake) Create(_ context.Context, configuracao *smtp.Configuracao) error {
	if r.configuracao != nil {
		return smtp.ErrorConfiguracaoSMTPAlreadyExists
	}
	copia := *configuracao
	r.configuracao = &copia
	return nil
}
func (r *smtpRepositoryControllerFake) Update(_ context.Context, configuracao *smtp.Configuracao) error {
	copia := *configuracao
	r.configuracao = &copia
	return nil
}

type smtpLoggerFake struct{}

func (smtpLoggerFake) Info(string, ...any)    {}
func (smtpLoggerFake) Warn(string, ...any)    {}
func (smtpLoggerFake) Error(string, ...any)   {}
func (smtpLoggerFake) Debug(string, ...any)   {}
func (smtpLoggerFake) With(...any) log.Logger { return smtpLoggerFake{} }
func (smtpLoggerFake) Fatal(string, ...any)   {}
func (smtpLoggerFake) Sync() error            { return nil }

func smtpControllerTeste(repo *smtpRepositoryControllerFake) *SMTPController {
	return NewSMTPController(smtp.NewService(repo, smtpLoggerFake{}))
}

func TestSMTPResponseNuncaSerializaSenha(t *testing.T) {
	response := toSMTPResponse(&smtp.Configuracao{Id: uuid.New(), Host: "smtp.example.com", Porta: 587, Usuario: "usuario", Senha: "segredo", Remetente: "no-reply@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	body, err := json.Marshal(response)
	if err != nil || strings.Contains(string(body), "segredo") || bytes.Contains(body, []byte("smtp_password")) {
		t.Fatalf("resposta expôs senha: %s / %v", body, err)
	}
}

func TestUpdateSMTPRejeitaCampoNuloEVazio(t *testing.T) {
	for _, body := range []string{`{}`, `{"smtp_password":null}`} {
		gin.SetMode(gin.TestMode)
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/admin/config-smtp", strings.NewReader(body))
		ctx.Request.Header.Set("Content-Type", "application/json")
		(&SMTPController{}).Update(ctx)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("body %s: status = %d", body, recorder.Code)
		}
	}
}

func TestSMTPControllerCriaConsultaEAtualizaSemExporSenha(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &smtpRepositoryControllerFake{}
	controller := smtpControllerTeste(repo)
	createRecorder := httptest.NewRecorder()
	createCtx, _ := gin.CreateTestContext(createRecorder)
	createCtx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/config-smtp", strings.NewReader(`{"smtp_host":"smtp.example.com","smtp_port":587,"smtp_username":"usuario","smtp_password":"segredo","smtp_from":"no-reply@example.com"}`))
	createCtx.Request.Header.Set("Content-Type", "application/json")
	controller.Create(createCtx)
	if createRecorder.Code != http.StatusCreated || strings.Contains(createRecorder.Body.String(), "segredo") {
		t.Fatalf("criação = %d/%s", createRecorder.Code, createRecorder.Body.String())
	}

	host := "smtp.novo.example.com"
	patchRecorder := httptest.NewRecorder()
	patchCtx, _ := gin.CreateTestContext(patchRecorder)
	patchCtx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/admin/config-smtp", strings.NewReader(`{"smtp_host":"`+host+`"}`))
	patchCtx.Request.Header.Set("Content-Type", "application/json")
	controller.Update(patchCtx)
	if patchRecorder.Code != http.StatusOK || repo.configuracao.Senha != "segredo" {
		t.Fatalf("atualização = %d; config = %+v", patchRecorder.Code, repo.configuracao)
	}

	getRecorder := httptest.NewRecorder()
	getCtx, _ := gin.CreateTestContext(getRecorder)
	getCtx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/config-smtp", nil)
	controller.Get(getCtx)
	if getRecorder.Code != http.StatusOK || strings.Contains(getRecorder.Body.String(), "segredo") {
		t.Fatalf("consulta = %d/%s", getRecorder.Code, getRecorder.Body.String())
	}
}
