package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"main.go/dto"
	"main.go/service"
)

// Cria um novo colaboradores
func ColaboradorCreate(c *gin.Context) {
	var colaborador dto.ColabotadoresRequestDTO

	c.Bind(&colaborador)

	colaboradorCreated, err := service.CreateColaboradores(colaborador)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"colaborador": colaboradorCreated,
	})
}

// Pega todos os usuários
func ColaboradoresGetAll(c *gin.Context) {
	colaboradores, err := service.GetAllColaboradores()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"colaboradores": colaboradores,
	})
}

// Pega um usuário pelo ID
func ColaboradorGetById(c *gin.Context) {
	idParam := c.Param("id")

	colaborador, err := service.GetColaboradorById(idParam)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"colaborador": colaborador,
	})
}

// Atualiza um usuário pelo ID
func ColaboradorUpdate(c *gin.Context) {
	idParam := c.Param("id")
	var body dto.ColabotadoresRequestDTO

	c.Bind(&body)

	colaboradorUpdated, err := service.UpdateColaborador(idParam, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"colaborador": colaboradorUpdated,
	})
}

// Deleta um usuário pelo ID
func ColaboradorDelete(c *gin.Context) {
	idParam := c.Param("id")
	err := service.DeleteColaboradorById(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Colaborador deleted successfully",
	})
}
