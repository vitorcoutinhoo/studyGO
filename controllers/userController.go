package controllers

import (
	"net/http"
	"strconv"

	"main.go/database"
	"main.go/models"

	"github.com/gin-gonic/gin"
)

// struct parametro para os endpoinds de post
type UserRequest struct {
	Email     string `json:"email"`
	SenhaHash string `json:"senhaHash"`
}

// UserCreate godoc
// @Summary Cria um novo usuário
// @Description Cria um novo usuário com email e senha hash
// @Tags Usuários
// @Accept json
// @Produce json
// @Param user body controllers.UserRequest true "Dados do usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func UserCreate(c *gin.Context) {
	var req UserRequest
	c.Bind(&req)

	userCreated := models.User{
		Email:     req.Email,
		SenhaHash: req.SenhaHash,
		Role:      "USER",
	}

	result := database.DB.Create(&userCreated)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userCreated})
}

// UserGetAll godoc
// @Summary Lista todos os usuários
// @Description Retorna todos os usuários cadastrados
// @Tags Usuários
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func UserGetAll(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// UserGetById godoc
// @Summary Busca um usuário pelo ID
// @Description Retorna um usuário específico com base no seu ID
// @Tags Usuários
// @Produce json
// @Param id path uint true "ID do usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/{id} [get]
func UserGetById(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	result := database.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UserUpdate godoc
// @Summary Atualiza um usuário pelo ID
// @Description Atualiza os dados (email e senha) de um usuário específico
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path uint true "ID do usuário"
// @Param user body controllers.UserRequest true "Novos dados do usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/{id} [put]
func UserUpdate(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req UserRequest
	c.Bind(&req)

	database.DB.Model(&user).Updates(models.User{
		Email:     req.Email,
		SenhaHash: req.SenhaHash,
	})

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UserDelete godoc
// @Summary Deleta um usuário pelo ID
// @Description Remove um usuário do banco de dados
// @Tags Usuários
// @Produce json
// @Param id path uint true "ID do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/{id} [delete]
func UserDelete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
