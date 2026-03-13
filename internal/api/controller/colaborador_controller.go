package controller

import (
	"fmt"
	"io"
	"net/http"
	"plantao/internal/api/dto"
	"plantao/internal/domain/colaborador"
	"plantao/internal/infra/config"
	"plantao/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// ColaboradorController é responsável por lidar com as requisições relacionadas aos colaboradores.
type ColaboradorController struct {
	service   *colaborador.ColaboradorService
	urlServer string
}

// NewColaboradorController cria uma nova instância de ColaboradorController.
func NewColaboradorController(service *colaborador.ColaboradorService, cfg *config.Config) *ColaboradorController {
	return &ColaboradorController{
		service:   service,
		urlServer: "http://" + cfg.Server.Host + ":" + cfg.Server.Port,
	}
} // Fim NewColaboradorController

// CreateColaborador lida com a criação de um novo colaborador.
func (c *ColaboradorController) CreateColaborador(ctx *gin.Context) {
	var req dto.CreateColaboradorRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var file io.ReadSeeker
	filename := ""

	uploadedFile, header, err := ctx.Request.FormFile("foto")

	// Se a foto foi enviada
	if err == nil {
		defer uploadedFile.Close()
		file = uploadedFile
		filename = header.Filename
	}

	col, err := createColaboradorDtoToDomain(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col.Foto = filename
	result, err := c.service.CreateColaborador(ctx, col, file)
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

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var file io.ReadSeeker
	filename := ""

	uploadedFile, header, err := ctx.Request.FormFile("foto")

	// Se a foto foi enviada
	if err == nil {
		defer uploadedFile.Close()
		file = uploadedFile
		filename = header.Filename
	}

	col, err := updateColaboradorDtoToDomain(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	col.Foto = filename
	if err := c.service.UpdateColaborador(ctx, col, id, file); err != nil {
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

	if resp.Foto != "" {
		resp.Foto = c.urlServer + resp.Foto
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

		if resp.Foto != "" {
			resp.Foto = c.urlServer + resp.Foto
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
		Foto:             "",
		Status:           ativo,
		AtivoPlantao:     ativoPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
		Cargo:            cargo,
	}, nil
}

func updateColaboradorDtoToDomain(r *dto.UpdateColaboradorRequest) (*colaborador.Colaborador, error) {
	var dataAdmissao *time.Time
	var err error

	if r.DataAdmissao != nil {
		dataAdmissao, err = utils.ParseBrToUsDate(r.DataAdmissao)
		if err != nil {
			return nil, err
		}
	} else {
		dataAdmissao = nil
	}

	var dataDesligamento *time.Time
	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento)
		if err != nil {
			return nil, err
		}
		dataDesligamento = dataTemp
	}

	var ativo *colaborador.StatusColaborador
	if r.Status != nil {
		status, err := colaborador.ParseStatusColaborador(*r.Status)
		if err != nil {
			return nil, err
		}
		ativo = &status
	} else {
		zero := colaborador.StatusColaborador(0)
		ativo = &zero
	}

	var ativoPlantao *colaborador.StatusColaborador
	if r.AtivoPlantao != nil {
		status, err := colaborador.ParseStatusColaborador(*r.AtivoPlantao)
		if err != nil {
			return nil, err
		}
		ativoPlantao = &status
	} else {
		zero := colaborador.StatusColaborador(0)
		ativoPlantao = &zero
	}

	var cargo *colaborador.CargoColaborador
	if r.Cargo != nil {
		c, err := colaborador.ParseCargoColaborador(*r.Cargo)
		if err != nil {
			return nil, err
		}
		cargo = &c
	} else {
		zero := colaborador.CargoColaborador("")
		cargo = &zero
	}

	var setor *colaborador.SetorColaborador
	if r.Setor != nil {
		s, err := colaborador.ParseSetorColaborador(*r.Setor)
		if err != nil {
			return nil, err
		}
		setor = &s
	} else {
		zero := colaborador.SetorColaborador("")
		setor = &zero
	}

	return &colaborador.Colaborador{
		Nome:             *r.Nome,
		Email:            *r.Email,
		Telefone:         *r.Telefone,
		Cargo:            *cargo,
		Setor:            *setor,
		Foto:             "",
		Status:           *ativo,
		AtivoPlantao:     *ativoPlantao,
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
