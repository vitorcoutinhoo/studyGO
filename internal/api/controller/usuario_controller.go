package controller

import (
	"net/http"
	"plantao/internal/api/dto"
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

func (c *UsuarioController) CreateUsuario(ctx *gin.Context) {
	idColaborador := ctx.Param("id_colaborador")

	var req dto.UsuarioRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuario, err := c.service.CreateUsuario(ctx, &req, idColaborador)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, usuario)
}

func (c *UsuarioController) UpdateUsuario(ctx *gin.Context) {
	idUsuario := ctx.Param("id_usuario")

	var req dto.UsuarioRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.UpdateUsuario(ctx, &req, idUsuario)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UsuarioController) DisableUsuario(ctx *gin.Context) {
	idUsuario := ctx.Param("id_usuario")

	err := c.service.DisableUsuario(ctx, idUsuario)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) GetUsuarioById(ctx *gin.Context) {
	idUsuario := ctx.Param("id_usuario")

	usuario, err := c.service.GetUsuarioById(ctx, idUsuario)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, usuario)
}
