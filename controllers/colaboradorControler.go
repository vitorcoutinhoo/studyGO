package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main.go/types"
)

type ColaboradorController struct {
	repository types.ColaboradorRepository
}

func NewColaboradorController(repository types.ColaboradorRepository) *ColaboradorController {
	return &ColaboradorController{
		repository: repository,
	}
}

func (c *ColaboradorController) RegisterRoutes(r *gin.Engine) {
	r.POST("/colaborador", c.colaboradorRegister)
	r.GET("/colaborador", c.colaboradorGet)
	r.GET("/colaborador/:id", c.colaboradorGetById)
	r.PUT("/colaborador/:id", c.colaboradorUpdate)
	r.DELETE("/colaborador/:id", c.colaboradorDelete)
}

func (c *ColaboradorController) colaboradorRegister(g *gin.Context) {
	var colaborador types.ColaboradorRequest

	g.Bind(&colaborador)

	colaboradorSaved, err := c.repository.CreateColaborador(colaborador)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	g.JSON(http.StatusOK, gin.H{
		"user": colaboradorSaved,
	})
}

func (c *ColaboradorController) colaboradorGet(g *gin.Context) {
	var colaborador types.ColaboradorRequest

	g.Bind(&colaborador)

	colaboradores, err := c.repository.GetColaboradores()

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	g.JSON(http.StatusOK, gin.H{
		"users": colaboradores,
	})
}

func (c *ColaboradorController) colaboradorGetById(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	colaborador, err := c.repository.GetColaboradoresById(idConverted)

	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"user": colaborador,
	})
}

func (c *ColaboradorController) colaboradorUpdate(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	var colaborador types.ColaboradorRequest
	g.Bind(&colaborador)

	err := c.repository.UpdateColaborador(idConverted, colaborador)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Colaborador atualizado com sucesso",
	})
}

func (c *ColaboradorController) colaboradorDelete(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	err := c.repository.DeleteColaboradorById(idConverted)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Colaborador deletado com sucesso",
	})
}
