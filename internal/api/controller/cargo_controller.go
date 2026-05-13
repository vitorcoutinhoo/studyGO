package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/domain/cargo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CargoController struct {
	service *cargo.CargoService
}

func NewCargoController(service *cargo.CargoService) *CargoController {
	return &CargoController{service: service}
}

func (c *CargoController) GetAll(ctx *gin.Context) {
	cargos, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	type response struct {
		Id   string `json:"id"`
		Nome string `json:"nome"`
	}

	resp := make([]response, 0, len(cargos))
	for _, ca := range cargos {
		resp = append(resp, response{Id: ca.Id.String(), Nome: ca.Nome})
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *CargoController) Create(ctx *gin.Context) {
	var req struct {
		Nome string `json:"nome" binding:"required,min=2,max=100"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	ca, err := c.service.Create(ctx.Request.Context(), req.Nome)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": ca.Id.String(), "nome": ca.Nome})
}

func (c *CargoController) Delete(ctx *gin.Context) {
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
