package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/shared"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PlantaoController struct {
	service *plantao.PlantaoService
}

func NewPlantaoController(service *plantao.PlantaoService) *PlantaoController {
	return &PlantaoController{
		service: service,
	}
}

func (p *PlantaoController) CreatePlantao(ctx *gin.Context) {
	var req dto.CreatePlantaoRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	inicio, err := time.Parse("2006-01-02", req.Periodo.Inicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data de início inválida"})
		return
	}

	fim, err := time.Parse("2006-01-02", req.Periodo.Fim)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data de fim inválida"})
		return
	}

	periodo := &shared.Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	plantaoCriado, err := p.service.CreatePlantao(ctx.Request.Context(), req.ColaboradorId, periodo)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toPlantaoResponse(plantaoCriado))
}

func (p *PlantaoController) UpdateStatusPlantao(ctx *gin.Context) {
	var req dto.UpdateStatusPlantaoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	status, err := strconv.Atoi(req.NewStatus)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "formato de status inválido"})
		return
	}

	plantaoId := ctx.Param("id")
	_, err = p.service.UpdatePlantaoStatus(ctx.Request.Context(), plantaoId, plantao.StatusPlantao(status), req.Observacoes)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (p *PlantaoController) GetPlantaoById(ctx *gin.Context) {
	plantaoId := ctx.Param("id")
	pl, err := p.service.GetPlantaoById(ctx.Request.Context(), plantaoId)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toPlantaoResponse(pl))
}

func (p *PlantaoController) GetPlantoes(ctx *gin.Context) {
	plantoes, err := p.service.GetPlantoes(ctx.Request.Context(), &plantao.Filtro{})
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	var response []dto.CreatePlantaoResponse
	for _, pl := range plantoes {
		response = append(response, *toPlantaoResponse(&pl))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) GetPlantoesByColaboradorId(ctx *gin.Context) {
	colaboradorId := ctx.Param("colaborador_id")
	plantoes, err := p.service.GetPlantoesByColaboradorId(ctx.Request.Context(), colaboradorId)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	var response []dto.CreatePlantaoResponse
	for _, pl := range plantoes {
		response = append(response, *toPlantaoResponse(&pl))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) GetPlantoesByPeriodo(ctx *gin.Context) {
	inicioStr := ctx.Param("start_date")
	fimStr := ctx.Param("end_date")

	if inicioStr == "" || fimStr == "" {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "start_date e end_date são obrigatórios (YYYY-MM-DD)"})
		return
	}

	inicio, err := time.Parse("2006-01-02", inicioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "start_date inválida"})
		return
	}

	fim, err := time.Parse("2006-01-02", fimStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "end_date inválida"})
		return
	}

	periodo := shared.Periodo{Inicio: inicio, Fim: fim}
	plantoes, err := p.service.GetPlantoesByPeriodo(ctx.Request.Context(), &periodo)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	var response []dto.CreatePlantaoResponse
	for _, pl := range plantoes {
		response = append(response, *toPlantaoResponse(&pl))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) GetPlantoesByStatus(ctx *gin.Context) {
	statusStr := ctx.Param("status")

	statusNum, err := strconv.Atoi(statusStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "formato de status inválido"})
		return
	}

	plantoes, err := p.service.GetPlantoesByStatus(ctx.Request.Context(), plantao.StatusPlantao(statusNum))
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	var response []dto.CreatePlantaoResponse
	for _, pl := range plantoes {
		response = append(response, *toPlantaoResponse(&pl))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) DeletePlantao(ctx *gin.Context) {
	plantaoId := ctx.Param("id")
	if err := p.service.DeletePlantao(ctx.Request.Context(), plantaoId); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func toPlantaoResponse(pl *plantao.Plantao) *dto.CreatePlantaoResponse {
	return &dto.CreatePlantaoResponse{
		Id:            pl.Id,
		ColaboradorId: pl.ColaboradorId,
		Periodo:       *pl.Periodo,
		Status:        pl.Status,
		ValorTotal:    pl.ValorTotal,
		Observacoes:   pl.Observacoes,
	}
}
