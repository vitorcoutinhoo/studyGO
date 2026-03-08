package controller

import (
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/comunicacao"

	"github.com/gin-gonic/gin"
)

type ModeloComunicacaoController struct {
	service *comunicacao.ModeloComunicacaoService
}

func NewModeloComunicacaoController(service *comunicacao.ModeloComunicacaoService) *ModeloComunicacaoController {
	return &ModeloComunicacaoController{
		service: service,
	}
}

func (c *ModeloComunicacaoController) CreateModeloComunicacao(ctx *gin.Context) {
	var req dto.ModeloComunicacaoRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelo, err := c.service.CreateModeloComunicacao(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, modelo)
}

func (c *ModeloComunicacaoController) UpdateModeloComunicacao(ctx *gin.Context) {
	idModelo := ctx.Param("id_modelo")

	var req dto.ModeloComunicacaoUpdateRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.UpdateModeloComunicacao(ctx, idModelo, &req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ModeloComunicacaoController) DisableModeloComunicacao(ctx *gin.Context) {
	idModelo := ctx.Param("id_modelo")

	err := c.service.DisableModeloComunicacao(ctx, idModelo)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ModeloComunicacaoController) GetModeloComunicacaoById(ctx *gin.Context) {
	idModelo := ctx.Param("id_modelo")

	modelo, err := c.service.GetModeloComunicacaoById(ctx, idModelo)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, modelo)
}

func (c *ModeloComunicacaoController) GetAllModelosComunicacao(ctx *gin.Context) {
	modelos, err := c.service.GetAllModelosComunicacao(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, modelos)
}
