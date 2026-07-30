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

func (c *ValorDiaController) UpdateValorVigente(ctx *gin.Context) {
	tipoDia := ctx.Param("tipo_dia")
	if tipoDia == "" {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "tipo_dia é obrigatório na URL"})
		return
	}

	var req dto.UpdateValorDiaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}
	if !req.HasFields() {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: financeiro.ErrorAtualizacaoVazia.Error()})
		return
	}

	atualizacao := &financeiro.AtualizacaoValorDia{
		DescricaoInformada:   req.Descricao.Set,
		Descricao:            req.Descricao.Value,
		VigenciaFimInformada: req.VigenciaFim.Set,
	}

	if req.Valor.Set {
		if req.Valor.Value == nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "valor não pode ser null"})
			return
		}
		atualizacao.Valor = req.Valor.Value
	}
	if req.VigenciaInicio.Set {
		if req.VigenciaInicio.Value == nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "vigencia_inicio não pode ser null"})
			return
		}
		data, err := parseValorDiaDate(*req.VigenciaInicio.Value)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "vigencia_inicio inválida, use o formato YYYY-MM-DD"})
			return
		}
		atualizacao.VigenciaInicio = &data
	}
	if req.VigenciaFim.Set && req.VigenciaFim.Value != nil {
		data, err := parseValorDiaDate(*req.VigenciaFim.Value)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "vigencia_fim inválida, use o formato YYYY-MM-DD ou null"})
			return
		}
		atualizacao.VigenciaFim = &data
	}

	valor, err := c.service.UpdateVigente(ctx.Request.Context(), financeiro.TipoDia(tipoDia), atualizacao)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.UpdateValorDiaResponse{
		Id:             valor.Id.String(),
		TipoDia:        string(valor.TipoDia),
		Valor:          valor.Valor,
		Descricao:      valor.Descricao,
		VigenciaInicio: valor.VigenciaInicio,
		VigenciaFim:    valor.VigenciaFim,
		UpdatedAt:      valor.UpdatedAt,
	})
}

func parseValorDiaDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
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
