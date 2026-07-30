package apihttp

import (
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
