package controller

import (
	"fmt"
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/colaborador"
	"plantao/internal/utils"
	"time"

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

	col, err := createColaboradorDtoToDomain(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.CreateColaborador(ctx, col)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := colaboradorToResponse(result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
} // Fim CreateColaborador

// UpdateColaborador lida com a atualização de um colaborador existente.
func (c *ColaboradorController) UpdateColaborador(ctx *gin.Context) {
	var req dto.UpdateColaboradorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := updateColaboradorDtoToDomain(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	if err := c.service.UpdateColaborador(ctx, col, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
} // Fim UpdateColaborador

// GetColaboradorById lida com a obtenção de um colaborador por ID.
func (c *ColaboradorController) GetColaboradorById(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.GetColaboradorById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := colaboradorToResponse(result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
} // Fim GetColaboradorById

// GetColaboradoresByFilter lida com a obtenção de colaboradores com base em filtros opcionais.
func (c *ColaboradorController) GetColaboradoresByFilter(ctx *gin.Context) {
	var filter dto.GetColaboradoresByFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := c.service.GetColaboradorByFilter(ctx, filterDtoToFilterDomain(filter))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.ColaboradorResponse, 0, len(results))
	for i := range results {
		resp, err := colaboradorToResponse(&results[i])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		responses = append(responses, *resp)
	}

	ctx.JSON(http.StatusOK, responses)
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

func createColaboradorDtoToDomain(r *dto.CreateColaboradorRequest) (*colaborador.Colaborador, error) {
	if r.Status == "" {
		r.Status = "ativo"
	}

	if r.AtivoPlantao == "" {
		r.AtivoPlantao = "ativo"
	}

	dataAdmissao, err := utils.ParseBrToUsDate(&r.DataAdmissao)
	if err != nil {
		return nil, err
	}

	var dataDesligamento *time.Time
	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento)
		if err != nil {
			return nil, err
		}
		dataDesligamento = dataTemp
	}

	ativo, err := colaborador.ParseStatusColaborador(r.Status)
	if err != nil {
		return nil, err
	}

	ativoPlantao, err := colaborador.ParseStatusColaborador(r.AtivoPlantao)
	if err != nil {
		return nil, err
	}

	cargo, err := colaborador.ParseCargoColaborador(r.Cargo)
	if err != nil {
		return nil, err
	}

	setor, err := colaborador.ParseSetorColaborador(r.Setor)
	if err != nil {
		return nil, err
	}

	return &colaborador.Colaborador{
		Nome:             r.Nome,
		Email:            r.Email,
		Telefone:         r.Telefone,
		Setor:            setor,
		Foto:             r.Foto,
		Status:           ativo,
		AtivoPlantao:     ativoPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
		Cargo:            cargo,
	}, nil
}

func updateColaboradorDtoToDomain(r *dto.UpdateColaboradorRequest) (*colaborador.Colaborador, error) {
	dataAdmissao, err := utils.ParseBrToUsDate(r.DataAdmissao)
	if err != nil {
		return nil, err
	}

	var dataDesligamento *time.Time
	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento)
		if err != nil {
			return nil, err
		}
		dataDesligamento = dataTemp
	}

	ativo, err := colaborador.ParseStatusColaborador(*r.Status)
	if err != nil {
		return nil, err
	}

	ativoPlantao, err := colaborador.ParseStatusColaborador(*r.AtivoPlantao)
	if err != nil {
		return nil, err
	}

	cargo, err := colaborador.ParseCargoColaborador(*r.Cargo)
	if err != nil {
		return nil, err
	}

	setor, err := colaborador.ParseSetorColaborador(*r.Setor)
	if err != nil {
		return nil, err
	}

	return &colaborador.Colaborador{
		Nome:             *r.Nome,
		Email:            *r.Email,
		Telefone:         *r.Telefone,
		Cargo:            cargo,
		Setor:            setor,
		Foto:             *r.Foto,
		Status:           ativo,
		AtivoPlantao:     ativoPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
}

func colaboradorToResponse(c *colaborador.Colaborador) (*dto.ColaboradorResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("colaborador vazio ou nulo")
	}

	status, err := colaborador.StatusColaboradorString(c.Status)
	if err != nil {
		return nil, err
	}

	statusPlantao, err := colaborador.StatusColaboradorString(c.AtivoPlantao)
	if err != nil {
		return nil, err
	}

	dataAdmissao, err := utils.ParseUsToBrDate(c.DataAdmissao)
	if err != nil {
		return nil, err
	}

	var dataDesligamento string
	if c.DataDesligamento != nil {
		dataTemp, err := utils.ParseUsToBrDate(c.DataDesligamento)
		if err != nil {
			return nil, err
		}
		dataDesligamento = dataTemp
	}

	return &dto.ColaboradorResponse{
		Id:               c.Id.String(),
		Nome:             c.Nome,
		Email:            c.Email,
		Telefone:         c.Telefone,
		Cargo:            string(c.Cargo),
		Setor:            string(c.Setor),
		Foto:             c.Foto,
		Status:           status,
		AtivoPlantao:     statusPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
}

func filterDtoToFilterDomain(filterReq dto.GetColaboradoresByFilterRequest) colaborador.ColaboradorFilter {
	return colaborador.ColaboradorFilter{
		Nome:         filterReq.Nome,
		Email:        filterReq.Email,
		Telefone:     filterReq.Telefone,
		Cargo:        filterReq.Cargo,
		Departamento: filterReq.Setor,
	}
}
