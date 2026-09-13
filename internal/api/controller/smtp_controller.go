package controller

import (
	"net/http"

	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/smtp"

	"github.com/gin-gonic/gin"
)

type SMTPController struct{ service *smtp.Service }

func NewSMTPController(service *smtp.Service) *SMTPController {
	return &SMTPController{service: service}
}

func (c *SMTPController) Get(ctx *gin.Context) {
	configuracao, err := c.service.Get(ctx.Request.Context())
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toSMTPResponse(configuracao))
}

func (c *SMTPController) Create(ctx *gin.Context) {
	var req dto.CreateSMTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apierr.RespondBinding(ctx, err)
		return
	}
	configuracao, err := c.service.Create(ctx.Request.Context(), req.Host, req.Porta, req.Usuario, req.Senha, req.Remetente)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toSMTPResponse(configuracao))
}

func (c *SMTPController) Update(ctx *gin.Context) {
	var req dto.UpdateSMTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apierr.RespondBinding(ctx, err)
		return
	}
	if !req.HasFields() || req.HasNullField() {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: smtp.ErrorAtualizacaoSMTPVazia.Error()})
		return
	}
	configuracao, err := c.service.Update(ctx.Request.Context(), &smtp.AtualizacaoConfiguracao{Host: req.Host.Value, Porta: req.Porta.Value, Usuario: req.Usuario.Value, Senha: req.Senha.Value, Remetente: req.Remetente.Value})
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toSMTPResponse(configuracao))
}

func toSMTPResponse(configuracao *smtp.Configuracao) dto.SMTPResponse {
	return dto.SMTPResponse{Id: configuracao.Id.String(), Host: configuracao.Host, Porta: configuracao.Porta, Usuario: configuracao.Usuario, Remetente: configuracao.Remetente, CreatedAt: configuracao.CreatedAt, UpdatedAt: configuracao.UpdatedAt}
}
