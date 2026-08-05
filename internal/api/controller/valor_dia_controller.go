package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/financeiro"

	"github.com/gin-gonic/gin"
)

type ValorDiaController struct {
	service *financeiro.ConfigValorDiaService
}

func NewValorDiaController(service *financeiro.ConfigValorDiaService) *ValorDiaController {
	return &ValorDiaController{service: service}
}

func (c *ValorDiaController) GetAll(ctx *gin.Context) {
	valores, err := c.service.GetAll(ctx.Request.Context())
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
	if req.HasDeprecatedVigenciaFields() {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "vigencia_inicio e vigencia_fim não são mais aceitos"})
		return
	}

	valor, err := c.service.SetValor(ctx.Request.Context(), financeiro.TipoDia(req.TipoDia), req.Valor, req.Descricao)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toValorDiaResponse(*valor))
}

func (c *ValorDiaController) UpdateValor(ctx *gin.Context) {
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
	if req.HasDeprecatedVigenciaFields() {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "vigencia_inicio e vigencia_fim não são mais aceitos"})
		return
	}
	if !req.HasFields() {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: financeiro.ErrorAtualizacaoVazia.Error()})
		return
	}

	atualizacao := &financeiro.AtualizacaoValorDia{
		DescricaoInformada: req.Descricao.Set,
		Descricao:          req.Descricao.Value,
	}

	if req.Valor.Set {
		if req.Valor.Value == nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "valor não pode ser null"})
			return
		}
		atualizacao.Valor = req.Valor.Value
	}
	valor, err := c.service.Update(ctx.Request.Context(), financeiro.TipoDia(tipoDia), atualizacao)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.UpdateValorDiaResponse{
		Id:        valor.Id.String(),
		TipoDia:   string(valor.TipoDia),
		Valor:     valor.Valor,
		Descricao: valor.Descricao,
		UpdatedAt: valor.UpdatedAt,
	})
}

func toValorDiaResponse(v financeiro.ValorDia) dto.ValorDiaResponse {
	return dto.ValorDiaResponse{
		Id:        v.Id.String(),
		TipoDia:   string(v.TipoDia),
		Valor:     v.Valor,
		Descricao: v.Descricao,
	}
}
