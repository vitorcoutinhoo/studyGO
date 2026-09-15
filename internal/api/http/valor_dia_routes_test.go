package apihttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"plantao/internal/api/controller"
	"plantao/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func TestSetupValorDiaRoutesRegistraPatchComESemTipo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setupValorDiaRoutes(router, &controller.ValorDiaController{}, &middleware.AuthMidware{})

	rotas := make(map[string]bool)
	for _, rota := range router.Routes() {
		rotas[rota.Method+" "+rota.Path] = true
	}
	for _, esperada := range []string{
		"PATCH /api/v1/admin/config-valores",
		"PATCH /api/v1/admin/config-valores/:tipo_dia",
	} {
		if !rotas[esperada] {
			t.Fatalf("rota não registrada: %s", esperada)
		}
	}
}

func TestRotasFinanceirasPermitemAdminEFinanceiro(t *testing.T) {
	for _, tt := range []struct {
		role   string
		status int
	}{
		{ADMIN_ROLE, http.StatusNoContent},
		{FINANCEIRO_ROLE, http.StatusNoContent},
		{GERENTE_ROLE, http.StatusForbidden},
		{COLABORADOR_ROLE, http.StatusForbidden},
	} {
		t.Run(tt.role, func(t *testing.T) {
			router := gin.New()
			router.Use(func(ctx *gin.Context) { ctx.Set("role", tt.role) })
			router.PATCH("/admin/config-valores/:tipo", middleware.RoleMidware(ADMIN_ROLE, FINANCEIRO_ROLE), func(ctx *gin.Context) {
				ctx.Status(http.StatusNoContent)
			})
			router.PATCH("/admin/feriados/:id/data", middleware.RoleMidware(ADMIN_ROLE, FINANCEIRO_ROLE), func(ctx *gin.Context) {
				ctx.Status(http.StatusNoContent)
			})

			for _, path := range []string{"/admin/config-valores/UTIL", "/admin/feriados/id/data"} {
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPatch, path, nil))
				if recorder.Code != tt.status {
					t.Fatalf("%s: status = %d, esperado %d", path, recorder.Code, tt.status)
				}
			}
		})
	}
}

func TestFinanceiroAcessaSomenteRotasGeraisMinimas(t *testing.T) {
	tests := []struct {
		nome     string
		roles    []string
		esperado int
	}{
		{"autoatendimento", []string{COLABORADOR_ROLE, GERENTE_ROLE, ADMIN_ROLE, FINANCEIRO_ROLE}, http.StatusNoContent},
		{"consulta cargos", []string{ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE, FINANCEIRO_ROLE}, http.StatusNoContent},
		{"consulta setores", []string{ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE, FINANCEIRO_ROLE}, http.StatusNoContent},
		{"colaboradores", []string{ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE}, http.StatusForbidden},
		{"plantões operacionais", []string{ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE}, http.StatusForbidden},
		{"administração", []string{ADMIN_ROLE}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			router := gin.New()
			router.Use(func(ctx *gin.Context) { ctx.Set("role", FINANCEIRO_ROLE) })
			router.GET("/recurso", middleware.RoleMidware(tt.roles...), func(ctx *gin.Context) {
				ctx.Status(http.StatusNoContent)
			})

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/recurso", nil))
			if recorder.Code != tt.esperado {
				t.Fatalf("status = %d, esperado %d", recorder.Code, tt.esperado)
			}
		})
	}
}

func TestSetupSMTPRoutesRegistraEndpointsAdministrativos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setupSMTPRoutes(router, &controller.SMTPController{}, &middleware.AuthMidware{})

	rotas := make(map[string]bool)
	for _, rota := range router.Routes() {
		rotas[rota.Method+" "+rota.Path] = true
	}
	for _, esperada := range []string{
		"GET /api/v1/admin/config-smtp",
		"POST /api/v1/admin/config-smtp",
		"PATCH /api/v1/admin/config-smtp",
	} {
		if !rotas[esperada] {
			t.Fatalf("rota não registrada: %s", esperada)
		}
	}
}
