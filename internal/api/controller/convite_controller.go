package controller

import (
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cl, err := c.colaboradorService.GetColaboradorById(ctx, req.IdColaborador)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cl == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "colaborador não encontrado"})
		return
	}

	if err := c.service.CreateConvite(ctx, req.IdColaborador, req.ColaboradorNome, req.ColaboradorEmail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *ConviteController) GetAllConvites(ctx *gin.Context) {
	convites, err := c.service.GetAllConvites(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cvt := make([]dto.ConviteResponseDTO, 0, len(convites))
	for i := range convites {
		expData, _ := utils.ParseUsToBrDate(&convites[i].ExpiraEm)
		createdAt, _ := utils.ParseUsToBrDate(&convites[i].CreatedAt)

		conviteDto := dto.ConviteResponseDTO{
			Token:         convites[i].Token.String(),
			IdColaborador: convites[i].IdColaborador.String(),
			ExpiraEm:      expData,
			Usado:         convites[i].Usado,
			CreatedAt:     createdAt,
		}

		cvt = append(cvt, conviteDto)
	}

	ctx.JSON(http.StatusOK, cvt)
}

func (c *ConviteController) DisableConvite(ctx *gin.Context) {
	token := ctx.Param("token")

	if err := c.service.DisableConvite(ctx, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
