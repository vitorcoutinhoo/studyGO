package middleware

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/security"

	"github.com/gin-gonic/gin"
)

type AuthMidware struct {
	jwtService     *security.JWTService
	usuarioService *usuario.UsuarioService
}

func NewAuthMidware(jwtService *security.JWTService, usuarioService *usuario.UsuarioService) *AuthMidware {
	return &AuthMidware{
		jwtService:     jwtService,
		usuarioService: usuarioService,
	}
}

func (a *AuthMidware) AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "É necessário autenticar-se para acessar este recurso."})
			return
		}

		claims, err := a.jwtService.ValidateToken(authCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "A sessão é inválida ou expirou. Faça login novamente."})
			return
		}

		err = a.usuarioService.ExistsUsuarioById(c.Request.Context(), claims.UserId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "A conta associada à sessão não está mais disponível."})
			return
		}

		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleMidware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, apierr.ErrorResponse{Code: "FORBIDDEN", Message: "Não foi possível identificar a permissão da sua conta."})
			return
		}

		roleStr := role.(string)

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, apierr.ErrorResponse{Code: "FORBIDDEN", Message: "Sua conta não possui permissão para realizar esta ação."})
	}
}
