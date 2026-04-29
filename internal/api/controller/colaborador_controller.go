package controller

import (
	"fmt"
	"io"
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/colaborador"
	"plantao/internal/infra/config"
	"plantao/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const maxFotoSize = 5 << 20 // 5 MB

type ColaboradorController struct {
	service   *colaborador.ColaboradorService
	urlServer string
}

func NewColaboradorController(service *colaborador.ColaboradorService, cfg *config.Config) *ColaboradorController {
	return &ColaboradorController{
		service:   service,
		urlServer: "http://" + cfg.Server.Host + ":" + cfg.Server.Port,
	}
}

func (c *ColaboradorController) CreateColaborador(ctx *gin.Context) {
	var req dto.CreateColaboradorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	var file io.ReadSeeker
	filename := ""

	uploadedFile, header, err := ctx.Request.FormFile("foto")
	if err == nil {
		if header.Size > maxFotoSize {
			ctx.JSON(http.StatusRequestEntityTooLarge, apierr.ErrorResponse{Code: "FILE_TOO_LARGE", Message: "foto não pode exceder 5 MB"})
			return
		}
		defer uploadedFile.Close()
		file = uploadedFile
		filename = header.Filename
	}

	col, err := createColaboradorDtoToDomain(&req)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	col.Foto = filename
	result, err := c.service.CreateColaborador(ctx, col, file)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	resp, err := colaboradorToResponse(result)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *ColaboradorController) UpdateColaborador(ctx *gin.Context) {
	var req dto.UpdateColaboradorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	var file io.ReadSeeker
	filename := ""

	uploadedFile, header, err := ctx.Request.FormFile("foto")
	if err == nil {
		if header.Size > maxFotoSize {
			ctx.JSON(http.StatusRequestEntityTooLarge, apierr.ErrorResponse{Code: "FILE_TOO_LARGE", Message: "foto não pode exceder 5 MB"})
			return
		}
		defer uploadedFile.Close()
		file = uploadedFile
		filename = header.Filename
	}

	col, err := updateColaboradorDtoToDomain(&req)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	id := ctx.Param("id")

	col.Foto = filename
	if err := c.service.UpdateColaborador(ctx, col, id, file); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ColaboradorController) GetColaboradorById(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.GetColaboradorById(ctx, id)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	resp, err := colaboradorToResponse(result)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	if resp.Foto != "" {
		resp.Foto = c.urlServer + resp.Foto
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *ColaboradorController) GetColaboradoresByFilter(ctx *gin.Context) {
	var filter dto.GetColaboradoresByFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	f, err := filterDtoToFilterDomain(filter)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	results, err := c.service.GetColaboradorByFilter(ctx, f)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	responses := make([]dto.ColaboradorResponse, 0, len(results))
	for i := range results {
		resp, err := colaboradorToResponse(&results[i])
		if err != nil {
			apierr.Respond(ctx, err)
			return
		}

		if resp.Foto != "" {
			resp.Foto = c.urlServer + resp.Foto
		}

		responses = append(responses, *resp)
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *ColaboradorController) UploadFotoColaborador(ctx *gin.Context) {
	id := ctx.Param("id")

	uploadedFile, header, err := ctx.Request.FormFile("foto")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "foto não enviada"})
		return
	}
	defer uploadedFile.Close()

	if header.Size > maxFotoSize {
		ctx.JSON(http.StatusRequestEntityTooLarge, apierr.ErrorResponse{Code: "FILE_TOO_LARGE", Message: "foto não pode exceder 5 MB"})
		return
	}

	if err := c.service.UpdateColaborador(ctx, &colaborador.Colaborador{Foto: header.Filename}, id, uploadedFile); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ColaboradorController) DisableColaborador(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DisableColaborador(ctx, id); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func createColaboradorDtoToDomain(r *dto.CreateColaboradorRequest) (*colaborador.Colaborador, error) {
	if r.Status == "" {
		r.Status = "ativo"
	}

	if r.AtivoPlantao == "" {
		r.AtivoPlantao = "ativo"
	}

	dataAdmissao, err := utils.ParseBrToUsDate(&r.DataAdmissao, nil)
	if err != nil {
		return nil, err
	}

	var dataDesligamento *time.Time
	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento, nil)
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

	return colaborador.NewColaborador(
		r.Nome,
		r.Email,
		r.Telefone,
		"",
		dataAdmissao,
		dataDesligamento,
		ativo,
		ativoPlantao,
		cargo,
		setor,
	)
}

func updateColaboradorDtoToDomain(r *dto.UpdateColaboradorRequest) (*colaborador.Colaborador, error) {
	var dataAdmissao *time.Time
	var err error

	if r.DataAdmissao != nil {
		dataAdmissao, err = utils.ParseBrToUsDate(r.DataAdmissao, nil)
		if err != nil {
			return nil, err
		}
	} else {
		dataAdmissao = nil
	}

	var dataDesligamento *time.Time
	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento, nil)
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
	if r.Cargo != nil && *r.Cargo != "" {
		c, err := colaborador.ParseCargoColaborador(*r.Cargo)
		if err != nil {
			return nil, err
		}
		cargo = &c
	}

	var setor *colaborador.SetorColaborador
	if r.Setor != nil && *r.Setor != "" {
		s, err := colaborador.ParseSetorColaborador(*r.Setor)
		if err != nil {
			return nil, err
		}
		setor = &s
	}

	var n string
	if r.Nome != nil {
		n = *r.Nome
	}

	var e string
	if r.Email != nil {
		e = *r.Email
	}

	var t string
	if r.Telefone != nil {
		t = *r.Telefone
	}

	c := &colaborador.Colaborador{
		Nome:             n,
		Email:            e,
		Telefone:         t,
		Foto:             "",
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}

	if ativo != nil {
		c.Status = *ativo
	}
	if ativoPlantao != nil {
		c.AtivoPlantao = *ativoPlantao
	}
	if cargo != nil {
		c.Cargo = *cargo
	}
	if setor != nil {
		c.Setor = *setor
	}

	return c, nil
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

	dataAdmissao, err := utils.ParseUsToBrDate(c.DataAdmissao, nil)
	if err != nil {
		return nil, err
	}

	var dataDesligamento string
	if c.DataDesligamento != nil {
		dataTemp, err := utils.ParseUsToBrDate(c.DataDesligamento, nil)
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

func filterDtoToFilterDomain(filterReq dto.GetColaboradoresByFilterRequest) (colaborador.ColaboradorFilter, error) {
	var data *time.Time
	var err error

	if filterReq.DataAdmissao != nil {
		data, err = utils.ParseBrToUsDate(filterReq.DataAdmissao, nil)
		if err != nil {
			return colaborador.ColaboradorFilter{}, err
		}
	} else {
		data = nil
	}

	return colaborador.ColaboradorFilter{
		Nome:         filterReq.Nome,
		Email:        filterReq.Email,
		Telefone:     filterReq.Telefone,
		Cargo:        filterReq.Cargo,
		DataAdmissao: data,
	}, nil
}
