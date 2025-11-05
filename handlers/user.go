package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Exemplo de rota de usuário
// @Description Retorna um usuário fictício
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string
// @Router /users [get]
func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":   1,
		"name": "Vítor Coutinho",
	})
}
