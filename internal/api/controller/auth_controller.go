package controller

import (
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/usuario"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *usuario.AuthService
}

func NewAuthController(authService *usuario.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.authService.Authenticate(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
