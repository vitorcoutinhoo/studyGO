package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/domain/setor"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SetorController struct {
	service *setor.SetorService
}

func NewSetorController(service *setor.SetorService) *SetorController {
	return &SetorController{service: service}
}

func (c *SetorController) GetAll(ctx *gin.Context) {
	setores, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	type response struct {
		Id   string `json:"id"`
		Nome string `json:"nome"`
	}

	resp := make([]response, 0, len(setores))
	for _, s := range setores {
		resp = append(resp, response{Id: s.Id.String(), Nome: s.Nome})
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *SetorController) Create(ctx *gin.Context) {
	var req struct {
		Nome string `json:"nome" binding:"required,min=2,max=100"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	s, err := c.service.Create(ctx.Request.Context(), req.Nome)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": s.Id.String(), "nome": s.Nome})
}

func (c *SetorController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "id inválido"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
