package controller

import (
	"net/http"
	"plantao/internal/domain/plantao"
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

type CreatePlantaoRequest struct {
	Periodo       plantao.Periodo `json:"periodo" binding:"required"`
	ColaboradorId string          `json:"colaborador_id" binding:"required"`
}

type UpdateStatusPlantaoRequest struct {
	NewStatus plantao.StatusPlantao `json:"new_status" binding:"required"`
}

type CreatePlantaoResponse struct {
	Id            string                `json:"id"`
	ColaboradorId string                `json:"colaborador_id"`
	Periodo       plantao.Periodo       `json:"periodo"`
	Status        plantao.StatusPlantao `json:"status"`
}

func (p *PlantaoController) CreatePlantao(ctx *gin.Context) {
	var req CreatePlantaoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plantao, err := p.service.CreatePlantao(ctx.Request.Context(), req.ColaboradorId, &req.Periodo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, toPlantaoResponse(plantao))
}

func (p *PlantaoController) UpdateStatusPlantao(ctx *gin.Context) {
	var req UpdateStatusPlantaoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plantaoId := ctx.Param("id")
	_, err := p.service.UpdatePlantaoStatus(ctx.Request.Context(), plantaoId, req.NewStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (p *PlantaoController) GetPlantaoById(ctx *gin.Context) {
	plantaoId := ctx.Param("id")
	plantao, err := p.service.GetPlantaoById(ctx.Request.Context(), plantaoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toPlantaoResponse(plantao))
}

func (p *PlantaoController) GetPlantoes(ctx *gin.Context) {
	plantoes, err := p.service.GetPlantoes(ctx.Request.Context(), nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []CreatePlantaoResponse
	for _, plantao := range plantoes {
		response = append(response, *toPlantaoResponse(&plantao))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) GetPlantoesByColaboradorId(ctx *gin.Context) {
	colaboradorId := ctx.Param("colaborador_id")
	plantoes, err := p.service.GetPlantoesByColaboradorId(ctx.Request.Context(), colaboradorId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []CreatePlantaoResponse
	for _, plantao := range plantoes {
		response = append(response, *toPlantaoResponse(&plantao))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) GetPlantoesByPeriodo(ctx *gin.Context) {
	iniciostr := ctx.Query("inicio")
	fimStr := ctx.Query("fim")

	inicio, err := time.Parse(time.RFC3339, iniciostr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid inicio format"})
		return
	}

	fim, err := time.Parse(time.RFC3339, fimStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fim format"})
		return
	}

	var periodo = plantao.Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	plantoes, err := p.service.GetPlantoesByPeriodo(ctx.Request.Context(), &periodo)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []CreatePlantaoResponse
	for _, plantao := range plantoes {
		response = append(response, *toPlantaoResponse(&plantao))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) GetPlantoesByStatus(ctx *gin.Context) {
	statusStr := ctx.Query("status")

	statusNum, err := strconv.Atoi(statusStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status format"})
		return
	}

	plantoes, err := p.service.GetPlantoesByStatus(ctx.Request.Context(), plantao.StatusPlantao(statusNum))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []CreatePlantaoResponse
	for _, plantao := range plantoes {
		response = append(response, *toPlantaoResponse(&plantao))
	}

	ctx.JSON(http.StatusOK, response)
}

func (p *PlantaoController) DeletePlantao(ctx *gin.Context) {
	plantaoId := ctx.Param("id")
	err := p.service.DeletePlantao(ctx.Request.Context(), plantaoId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func toPlantaoResponse(plantao *plantao.Plantao) *CreatePlantaoResponse {
	return &CreatePlantaoResponse{
		Id:            plantao.Id,
		ColaboradorId: plantao.ColaboradorId,
		Periodo:       *plantao.Periodo,
		Status:        plantao.Status,
	}
}
