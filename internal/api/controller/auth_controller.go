package controller

import (
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/config"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *usuario.AuthService
	expireTime  int64
}

func NewAuthController(authService *usuario.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		expireTime:  cfg.JWT.ExpireTime,
	}
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email é obrigatório"})
		return
	}

	if req.Senha == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "senha é obrigatória"})
		return
	}

	usuarioLogado, token, err := a.authService.Authenticate(ctx, *req.Email, *req.Senha)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	usuarioLogadoDto := dto.UsuarioResponseDTO{
		Id:            usuarioLogado.Id.String(),
		IdColaborador: usuarioLogado.IdColaborador.String(),
		Email:         usuarioLogado.Email,
		Role:          string(usuarioLogado.Role),
		Ativo:         convertStatusUsuario(usuarioLogado.Ativo),
	}

	a.setAccessTokenCookie(ctx, *token)

	ctx.JSON(http.StatusOK, usuarioLogadoDto)
}

func (a *AuthController) setAccessTokenCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",
		token,
		int(a.expireTime)*60,
		"/",
		"",
		false, // Secure: só HTTPS
		true,  // HttpOnly: JS não acessa
	)
}

func convertStatusUsuario(sts usuario.StatusUsuario) string {
	if sts == usuario.StatusAtivo {
		return "ativo"
	}

	return "inativo"
}
