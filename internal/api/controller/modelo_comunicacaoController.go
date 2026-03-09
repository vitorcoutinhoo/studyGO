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

	modelo, err := c.service.CreateModeloComunicacao(ctx, req.Nome, req.TipoComunicacao, req.Assunto, req.Corpo)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, modeloComunicacaoToResponse(*modelo))
}

func (c *ModeloComunicacaoController) UpdateModeloComunicacao(ctx *gin.Context) {
	idModelo := ctx.Param("id_modelo")

	var req dto.ModeloComunicacaoUpdateRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateModeloComunicacao(ctx, idModelo, req.Nome, req.TipoComunicacao, req.Assunto, req.Corpo, req.Ativo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ModeloComunicacaoController) DisableModeloComunicacao(ctx *gin.Context) {
	idModelo := ctx.Param("id_modelo")

	if err := c.service.DisableModeloComunicacao(ctx, idModelo); err != nil {
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

	ctx.JSON(http.StatusOK, modeloComunicacaoToResponse(*modelo))
}

func (c *ModeloComunicacaoController) GetAllModelosComunicacao(ctx *gin.Context) {
	modelos, err := c.service.GetAllModelosComunicacao(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.ModeloComunicacaoResponseDTO, 0, len(modelos))

	for _, m := range modelos {
		responses = append(responses, modeloComunicacaoToResponse(m))
	}

	ctx.JSON(http.StatusOK, responses)
}

func modeloComunicacaoToResponse(m comunicacao.Comunicacao) dto.ModeloComunicacaoResponseDTO {
	return dto.ModeloComunicacaoResponseDTO{
		Id:              m.Id.String(),
		Nome:            string(m.Nome),
		TipoComunicacao: string(m.TipoComunicacao),
		Assunto:         m.Assunto,
		Corpo:           m.Corpo,
		Ativo:           comunicacao.ParseStatusModeloComunicacaoString(m.Ativo),
	}
}
