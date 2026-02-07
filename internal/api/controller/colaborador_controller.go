package controller

import (
	"fmt"
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/colaborador"

	"github.com/gin-gonic/gin"
)

// ColaboradorController é responsável por lidar com as requisições relacionadas aos colaboradores.
type ColaboradorController struct {
	service *colaborador.ColaboradorService
}

// NewColaboradorController cria uma nova instância de ColaboradorController.
func NewColaboradorController(service *colaborador.ColaboradorService) *ColaboradorController {
	return &ColaboradorController{
		service: service,
	}
} // Fim NewColaboradorController

// CreateColaborador lida com a criação de um novo colaborador.
func (c *ColaboradorController) CreateColaborador(ctx *gin.Context) {
	var req dto.CreateColaboradorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	colaborador, err := c.service.CreateColaborador(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, colaborador)
} // Fim CreateColaborador

// UpdateColaborador lida com a atualização de um colaborador existente.
func (c *ColaboradorController) UpdateColaborador(ctx *gin.Context) {
	var req dto.UpdateColaboradorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	colaborador, err := c.service.UpdateColaborador(ctx, &req, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, colaborador)
} // Fim UpdateColaborador

// GetColaboradorById lida com a obtenção de um colaborador por ID.
func (c *ColaboradorController) GetColaboradorById(ctx *gin.Context) {
	id := ctx.Param("id")
	colaborador, err := c.service.GetColaboradorById(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, colaborador)
} // Fim GetColaboradorById

// GetColaboradoresByFilter lida com a obtenção de colaboradores com base em filtros opcionais.
func (c *ColaboradorController) GetColaboradoresByFilter(ctx *gin.Context) {
	var filter dto.GetColaboradoresByFilterRequest

	fmt.Println(filter)

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(filter)

	colaboradores, err := c.service.GetColaboradorByFilter(ctx, filter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, colaboradores)
} // Fim GetColaboradoresByFilter

// DisableColaborador lida com a desativação de um colaborador por ID.
func (c *ColaboradorController) DisableColaborador(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DisableColaborador(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
} // Fim DisableColaborador
