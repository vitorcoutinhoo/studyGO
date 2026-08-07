package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthenticationMiddlewareRejeitaRequisicaoSemToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protegida", (&AuthMidware{}).AuthenticationMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/protegida", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, esperado 401", recorder.Code)
	}
}

func TestRoleMidwarePermiteSomenteAdmin(t *testing.T) {
	tests := []struct {
		role     string
		status   int
		executou bool
	}{
		{"admin", http.StatusNoContent, true},
		{"gerente", http.StatusForbidden, false},
		{"colaborador", http.StatusForbidden, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			executou := false
			router.Use(func(c *gin.Context) {
				c.Set("role", tt.role)
			})
			router.PATCH("/admin/config-valores/UTIL", RoleMidware("admin"), func(c *gin.Context) {
				executou = true
				c.Status(http.StatusNoContent)
			})

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPatch, "/admin/config-valores/UTIL", nil))
			if recorder.Code != tt.status || executou != tt.executou {
				t.Fatalf("status/executou = %d/%v, esperado %d/%v", recorder.Code, executou, tt.status, tt.executou)
			}
		})
	}
}
