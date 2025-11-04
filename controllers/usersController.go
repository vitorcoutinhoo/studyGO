package controllers

import (
	"net/http"

	"main.go/initializers"
	"main.go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Cria um novo usuário
func UserCreate(c *gin.Context) {
	var user struct {
		Email     string
		SenhaHash string
	}

	c.Bind(&user)

	userCreated := models.User{
		ID:        uuid.New(),
		Email:     user.Email,
		SenhaHash: user.SenhaHash,
		Role:      "USER",
	}

	result := initializers.DB.Create(&userCreated)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userCreated,
	})
}

// Pega todos os usuários
func UserGetAll(c *gin.Context) {
	var posts []models.User
	initializers.DB.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"users": posts,
	})
}

// Pega um usuário pelo ID
func UserGetById(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
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

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ivalid ID",
		})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})

		return
	}

	var body struct {
		Email     string
		SenhaHash string
	}

	c.Bind(&body)

	initializers.DB.Model(&user).Updates(models.User{
		Email:     body.Email,
		SenhaHash: body.SenhaHash,
	})

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Deleta um usuário pelo ID
func UserDelete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})

		return
	}

	var user models.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})

		return
	}

	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}
