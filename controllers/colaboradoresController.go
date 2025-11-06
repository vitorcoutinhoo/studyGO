package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"main.go/db"
	"main.go/models"
)

// Cria um novo colaboradores
func ColaboradorCreate(c *gin.Context) {
	var colaborador struct {
		Nome         string
		Email        string
		Telefone     *string
		Cargo        *string
		Departamento *string
		FotoURL      string
		Ativo        string
		DataAdmissao time.Time
	}

	c.Bind(&colaborador)

	colaboradorCreated := models.Colaboradores{
		ID:           uuid.New(),
		Nome:         colaborador.Nome,
		Email:        colaborador.Email,
		Telefone:     colaborador.Telefone,
		Cargo:        colaborador.Cargo,
		Departamento: colaborador.Departamento,
		FotoURL:      colaborador.FotoURL,
		Ativo:        colaborador.Ativo,
		DataAdmissao: colaborador.DataAdmissao,
	}

	result := db.DB.Create(&colaboradorCreated)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"colaborador": colaboradorCreated,
	})
}

// Pega todos os usuários
func ColaboradoresGetAll(c *gin.Context) {
	var colaboradores []models.User
	db.DB.Find(&colaboradores)

	c.JSON(http.StatusOK, gin.H{
		"colaboradores": colaboradores,
	})
}
