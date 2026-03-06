package controller

import (
	"fmt"
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/security"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	jwtService     *security.JWTService
	usuarioService *usuario.UsuarioService
}

func NewAuthController(jwtService *security.JWTService, usuarioService *usuario.UsuarioService) *AuthController {
	return &AuthController{
		jwtService:     jwtService,
		usuarioService: usuarioService,
	}
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := a.usuarioService.GetUsuarioByEmail(ctx, req.Email)

	fmt.Println("Login ", u)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "email ou senha invalidos"})
		return
	}

	token, err := a.jwtService.GenerateToken(u.Id, u.Role)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
