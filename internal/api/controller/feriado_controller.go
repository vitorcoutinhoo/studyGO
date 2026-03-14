package controller

import (
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/financeiro"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FeriadoController struct {
	service *financeiro.FeriadoService
}

func NewFeriadoController(service *financeiro.FeriadoService) *FeriadoController {
	return &FeriadoController{service: service}
}

func (c *FeriadoController) GetFeriadosByAno(ctx *gin.Context) {
	anoStr := ctx.Query("ano")
	if anoStr == "" {
		anoStr = strconv.Itoa(time.Now().Year())
	}

	ano, err := strconv.Atoi(anoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ano inválido"})
		return
	}

	feriados, err := c.service.GetFeriadosByAno(ctx.Request.Context(), ano)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.FeriadoResponse
	for _, f := range feriados {
		response = append(response, dto.FeriadoResponse{
			Id:        f.Id.String(),
			Data:      f.Data,
			Nome:      f.Nome,
			Descricao: f.Descricao,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *FeriadoController) UpdateDataFeriado(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateDataFeriadoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novaData, err := time.Parse("2006-01-02", req.NovaData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nova_data inválida, use o formato YYYY-MM-DD"})
		return
	}

	feriado, err := c.service.UpdateDataFeriado(ctx.Request.Context(), id, novaData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.FeriadoResponse{
		Id:        feriado.Id.String(),
		Data:      feriado.Data,
		Nome:      feriado.Nome,
		Descricao: feriado.Descricao,
	})
}
