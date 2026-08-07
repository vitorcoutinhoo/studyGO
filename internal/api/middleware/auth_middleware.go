package middleware

import (
	"net/http"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token não fornecido"})
			return
		}

		claims, err := a.jwtService.ValidateToken(authCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		err = a.usuarioService.ExistsUsuarioById(c.Request.Context(), claims.UserId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "usuário não existe"})
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role não encontrada"})
			return
		}

		roleStr := role.(string)

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
	}
}
