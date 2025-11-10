package controllers

import (
	"net/http"

	"main.go/dto"
	"main.go/service"

	"github.com/gin-gonic/gin"
)

// Cria um novo usuário
func UserCreate(c *gin.Context) {
	var user dto.UserRequestDTO

	c.Bind(&user)

	userCreated, err := service.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userCreated,
	})
}

// Pega todos os usuários
func UserGetAll(c *gin.Context) {
	users, err := service.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// Pega um usuário pelo ID
func UserGetById(c *gin.Context) {
	idParam := c.Param("id")

	user, err := service.GetUserById(idParam)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Atualiza um usuário pelo ID
func UserUpdate(c *gin.Context) {
	idParam := c.Param("id")
	var body dto.UserRequestDTO

	c.Bind(&body)

	userUpdated, err := service.UpdateUser(idParam, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userUpdated,
	})
}

// Deleta um usuário pelo ID
func UserDelete(c *gin.Context) {
	idParam := c.Param("id")

	err := service.DeletUserById(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
