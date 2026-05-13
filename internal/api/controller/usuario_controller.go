package controller

import (
	"fmt"
	"net/http"
	"plantao/internal/api/apierr"
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

func (c *UsuarioController) CreateUsuarioByToken(ctx *gin.Context) {
	tokenStr := ctx.Query("token")
	if tokenStr == "" {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "token de criação ausente"})
		return
	}

	var req dto.CadastroByTokenRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	result, err := c.service.CreateUsuarioByToken(ctx, tokenStr, req.Senha)
	if err != nil {
		switch err {
		case convite.ErrorConviteNotFound, convite.ErrorConviteUsed, convite.ErrorConviteExpired:
			ctx.JSON(http.StatusGone, apierr.ErrorResponse{Code: "GONE", Message: err.Error()})
		default:
			apierr.Respond(ctx, err)
		}
		return
	}

	resp, _ := usuarioToResponse(result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *UsuarioController) UpdateUsuario(ctx *gin.Context) {
	idUsuarioRaw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "usuário não autenticado"})
		return
	}

	var req dto.UsuarioRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	if err := c.service.UpdateUsuario(ctx, req.Email, req.Senha, idUsuarioRaw.(string)); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UsuarioController) DeleteUsuario(ctx *gin.Context) {
	idUsuarioRaw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "usuário não autenticado"})
		return
	}

	if err := c.service.DeleteUsuario(ctx, idUsuarioRaw.(string)); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) GetUsuarioById(ctx *gin.Context) {
	idUsuarioRaw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "usuário não autenticado"})
		return
	}

	result, err := c.service.GetUsuarioById(ctx, idUsuarioRaw.(string))
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	resp, err := usuarioToResponse(result)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UsuarioController) GetAll(ctx *gin.Context) {
	results, err := c.service.GetAll(ctx)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	var usuariosDTO []dto.UsuarioResponseDTO
	for _, result := range *results {
		uDTO, err := usuarioToResponse(&result)
		if err != nil {
			apierr.Respond(ctx, err)
			return
		}
		usuariosDTO = append(usuariosDTO, *uDTO)
	}

	ctx.JSON(http.StatusOK, usuariosDTO)
}

func (c *UsuarioController) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin gerente colaborador"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	if err := c.service.UpdateRole(ctx, id, usuario.Role(req.Role)); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
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
