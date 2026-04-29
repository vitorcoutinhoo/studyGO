package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/financeiro"
	"time"

	"github.com/gin-gonic/gin"
)

type ValorDiaController struct {
	service *financeiro.ConfigValorDiaService
}

func NewValorDiaController(service *financeiro.ConfigValorDiaService) *ValorDiaController {
	return &ValorDiaController{service: service}
}

func (c *ValorDiaController) GetVigentes(ctx *gin.Context) {
	valores, err := c.service.GetVigentes(ctx.Request.Context())
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	var response []dto.ValorDiaResponse
	for _, v := range valores {
		response = append(response, toValorDiaResponse(v))
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ValorDiaController) SetValor(ctx *gin.Context) {
	var req dto.SetValorDiaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	vigenciaInicio, err := time.Parse("2006-01-02", req.VigenciaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "vigencia_inicio inválida, use o formato YYYY-MM-DD"})
		return
	}

	valor, err := c.service.SetValor(ctx.Request.Context(), financeiro.TipoDia(req.TipoDia), req.Valor, vigenciaInicio)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toValorDiaResponse(*valor))
}

func toValorDiaResponse(v financeiro.ValorDia) dto.ValorDiaResponse {
	return dto.ValorDiaResponse{
		Id:             v.Id.String(),
		TipoDia:        string(v.TipoDia),
		Valor:          v.Valor,
		VigenciaInicio: v.VigenciaInicio,
		VigenciaFim:    v.VigenciaFim,
	}
}
