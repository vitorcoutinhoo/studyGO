package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/api/dto"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/convite"
	"plantao/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConviteController struct {
	service            *convite.ConviteService
	colaboradorService *colaborador.ColaboradorService
}

func NewConviteController(service *convite.ConviteService, colaboradorService *colaborador.ColaboradorService) *ConviteController {
	return &ConviteController{
		service:            service,
		colaboradorService: colaboradorService,
	}
}

func (c *ConviteController) CreateConvite(ctx *gin.Context) {
	var req dto.ConviteRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apierr.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	cl, err := c.colaboradorService.GetColaboradorById(ctx, req.IdColaborador)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	if cl == nil {
		ctx.JSON(http.StatusNotFound, apierr.ErrorResponse{Code: "NOT_FOUND", Message: "colaborador não encontrado"})
		return
	}

	if err := c.service.CreateConvite(ctx, req.IdColaborador, req.ColaboradorNome, req.ColaboradorEmail); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *ConviteController) GetAllConvites(ctx *gin.Context) {
	convites, err := c.service.GetAllConvites(ctx)
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	cvt := make([]dto.ConviteResponseDTO, 0, len(convites))
	for i := range convites {
		expData, _ := utils.ParseUsToBrDate(&convites[i].ExpiraEm, nil)
		createdAt, _ := utils.ParseUsToBrDate(&convites[i].CreatedAt, nil)

		cvt = append(cvt, dto.ConviteResponseDTO{
			Token:         convites[i].Token.String(),
			IdColaborador: convites[i].IdColaborador.String(),
			ExpiraEm:      expData,
			Usado:         convites[i].Usado,
			CreatedAt:     createdAt,
		})
	}

	ctx.JSON(http.StatusOK, cvt)
}

func (c *ConviteController) DisableConvite(ctx *gin.Context) {
	token := ctx.Param("token")

	if err := c.service.DisableConvite(ctx, token); err != nil {
		apierr.Respond(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
