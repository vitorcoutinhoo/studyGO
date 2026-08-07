package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/shared"
	"plantao/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	inicio, err := parseDateTime(req.Periodo.Inicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data de início inválida, use YYYY-MM-DD ou YYYY-MM-DDTHH:MM:SSZ"})
		return
	}

	fim, err := parseDateTime(req.Periodo.Fim)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data de fim inválida, use YYYY-MM-DD ou YYYY-MM-DDTHH:MM:SSZ"})
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

func (p *PlantaoController) UpdatePlantao(ctx *gin.Context) {
	plantaoID := ctx.Param("id")
	if _, err := uuid.Parse(plantaoID); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "id do plantão inválido"})
		return
	}
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	var req dto.UpdatePlantaoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}
	if !req.HasFields() {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: plantao.ErrorAtualizacaoPlantaoVazia.Error()})
		return
	}

	atualizacao := &plantao.AtualizacaoPlantao{}
	if req.ColaboradorID.Set {
		if req.ColaboradorID.Value == nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "colaborador_id não pode ser null"})
			return
		}
		id, err := uuid.Parse(*req.ColaboradorID.Value)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "colaborador_id inválido"})
			return
		}
		colaboradorID := id.String()
		atualizacao.ColaboradorID = &colaboradorID
	}
	if req.DataInicio.Set {
		if req.DataInicio.Value == nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_inicio não pode ser null"})
			return
		}
		inicio, err := parseDateTime(*req.DataInicio.Value)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_inicio inválida, use YYYY-MM-DD ou RFC3339"})
			return
		}
		atualizacao.DataInicio = &inicio
	}
	if req.DataFim.Set {
		if req.DataFim.Value == nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_fim não pode ser null"})
			return
		}
		fim, err := parseDateTime(*req.DataFim.Value)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_fim inválida, use YYYY-MM-DD ou RFC3339"})
			return
		}
		atualizacao.DataFim = &fim
	}

	atualizado, err := p.service.EditarPlantao(ctx.Request.Context(), plantaoID, userID, atualizacao)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toPlantaoResponse(atualizado))
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
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}
	_, err = p.service.UpdatePlantaoStatus(ctx.Request.Context(), plantaoId, userID, plantao.StatusPlantao(status), req.Observacoes)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (p *PlantaoController) PagarPlantao(ctx *gin.Context) {
	var req dto.PagamentoPlantaoRequest
	if ctx.Request.ContentLength != 0 {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
			return
		}
	}

	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}
	if err := p.service.PagarPlantao(ctx.Request.Context(), ctx.Param("id"), userID, req.Observacoes); err != nil {
		apierr.Respond(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func authenticatedUserID(ctx *gin.Context) (string, bool) {
	raw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "usuário não autenticado"})
		return "", false
	}
	userID, ok := raw.(string)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, apierr.ErrorResponse{Code: "UNAUTHORIZED", Message: "identidade autenticada inválida"})
		return "", false
	}
	return userID, true
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

func (p *PlantaoController) GetRelatorio(ctx *gin.Context) {
	inicioStr := ctx.Query("data_inicio")
	fimStr := ctx.Query("data_fim")
	if inicioStr == "" || fimStr == "" {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_inicio e data_fim são obrigatórios (YYYY-MM-DD)"})
		return
	}

	inicio, err := time.Parse("2006-01-02", inicioStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_inicio inválida, use o formato YYYY-MM-DD"})
		return
	}
	fim, err := time.Parse("2006-01-02", fimStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "data_fim inválida, use o formato YYYY-MM-DD"})
		return
	}

	filtro := &plantao.RelatorioFiltro{DataInicio: inicio, DataFim: fim}
	if colaboradorID := ctx.Query("colaborador_id"); colaboradorID != "" {
		id, err := uuid.Parse(colaboradorID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "colaborador_id inválido"})
			return
		}
		filtro.ColaboradorID = id.String()
	}
	if statusStr := ctx.Query("status"); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: "status deve ser numérico"})
			return
		}
		statusPlantao := plantao.StatusPlantao(status)
		filtro.Status = &statusPlantao
	}

	items, err := p.service.GetRelatorio(ctx.Request.Context(), filtro)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	response := make([]dto.RelatorioPlantaoResponse, 0, len(items))
	for _, item := range items {
		response = append(response, dto.RelatorioPlantaoResponse{
			ColaboradorID:   item.ColaboradorID,
			PlantaoID:       item.PlantaoID,
			Status:          item.Status,
			DataInicio:      item.DataInicio,
			DataFim:         item.DataFim,
			Data:            item.Data.Format("2006-01-02"),
			NomeColaborador: item.NomeColaborador,
			ValorTotal:      item.ValorTotal,
			Valor:           item.Valor,
			Observacoes:     item.Observacoes,
		})
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

func parseDateTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02", s, utils.BrasilLocation())
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
