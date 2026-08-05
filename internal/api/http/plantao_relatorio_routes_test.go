package apihttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"plantao/internal/api/controller"
	"plantao/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func TestSetupPlantaoRoutesRegistraRelatorioProtegido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setupPlantaoRoutes(router, &controller.PlantaoController{}, &middleware.AuthMidware{})

	registrada := false
	for _, rota := range router.Routes() {
		if rota.Method == http.MethodGet && rota.Path == "/api/v1/plantoes/relatorio" {
			registrada = true
			break
		}
	}
	if !registrada {
		t.Fatal("rota GET /api/v1/plantoes/relatorio não registrada")
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status sem autenticação = %d, esperado 401", recorder.Code)
	}
}

func TestRelatorioPermiteAdminEGerenteERejeitaColaborador(t *testing.T) {
	tests := []struct {
		role   string
		status int
	}{
		{ADMIN_ROLE, http.StatusOK},
		{GERENTE_ROLE, http.StatusOK},
		{COLABORADOR_ROLE, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			router := gin.New()
			router.Use(func(ctx *gin.Context) {
				ctx.Set("role", tt.role)
			})
			router.GET("/relatorio", middleware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE), func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/relatorio", nil))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, esperado %d", recorder.Code, tt.status)
			}
		})
	}
}
