package midware

import (
	"net/http"
	"plantao/internal/infra/security"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMidware struct {
	jwtService *security.JWTService
}

func NewAuthMidware(jwtService *security.JWTService) *AuthMidware {
	return &AuthMidware{
		jwtService: jwtService,
	}
}

func (a *AuthMidware) AuthenticationMidware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token não fornecido"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := a.jwtService.ValidateToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role não encontrado"})
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
