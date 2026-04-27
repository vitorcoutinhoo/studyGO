package controller

import (
	"fmt"
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/convite"
	"plantao/internal/domain/usuario"

	"github.com/gin-gonic/gin"
)

type UsuarioController struct {
	service *usuario.UsuarioService
}

func NewUsuarioController(service *usuario.UsuarioService) *UsuarioController {
	return &UsuarioController{
		service: service,
	}
}

// No controller
func (c *UsuarioController) CreateUsuarioByToken(ctx *gin.Context) {
	tokenStr := ctx.Query("token") // ?token=<uuid>
	if tokenStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token de criação ausente"})
		return
	}

	var req dto.UsuarioRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.CreateUsuarioByToken(ctx, tokenStr, req.Email, req.Senha)
	if err != nil {
		status := http.StatusInternalServerError
		// Erros esperados retornam 400/410
		switch err {
		case convite.ErrorConviteNotFound, convite.ErrorConviteUsed, convite.ErrorConviteExpired:
			status = http.StatusGone // 410 — link expirado/usado
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	resp, _ := usuarioToResponse(result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *UsuarioController) UpdateUsuario(ctx *gin.Context) {
	idUsuarioRaw, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	var req dto.UsuarioRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateUsuario(ctx, req.Email, req.Senha, idUsuarioRaw.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UsuarioController) DeleteUsuario(ctx *gin.Context) {
	idUsuarioRaw, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	if err := c.service.DeleteUsuario(ctx, idUsuarioRaw.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) GetUsuarioById(ctx *gin.Context) {
	idUsuarioRaw, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	result, err := c.service.GetUsuarioById(ctx, idUsuarioRaw.(string))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := usuarioToResponse(result)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UsuarioController) GetAll(ctx *gin.Context) {
	results, err := c.service.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var usuariosDTO []dto.UsuarioResponseDTO

	for _, result := range *results {
		uDTO, err := usuarioToResponse(&result)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		usuariosDTO = append(usuariosDTO, *uDTO)
	}

	ctx.JSON(http.StatusOK, usuariosDTO)
}

func usuarioToResponse(u *usuario.Usuario) (*dto.UsuarioResponseDTO, error) {
	if u == nil {
		return nil, fmt.Errorf("usuário vazio ou nulo")
	}

	ativo, err := usuario.StatusUsuarioString(u.Ativo)

	if err != nil {
		return nil, err
	}

	return &dto.UsuarioResponseDTO{
		Id:            u.Id.String(),
		IdColaborador: u.IdColaborador.String(),
		Email:         u.Email,
		Role:          string(u.Role),
		Ativo:         ativo,
	}, nil
}
