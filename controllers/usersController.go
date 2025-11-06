package controllers

import (
	"net/http"

	"main.go/dto"
	"main.go/service"

	"github.com/gin-gonic/gin"
)

// Cria um novo usuário
func UserCreate(c *gin.Context) {
	var user dto.UserCreateDTO

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

/*// Pega um usuário pelo ID
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
	result := db.DB.First(&user, id)

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
	result := db.DB.First(&user, id)

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

	db.DB.Model(&user).Updates(models.User{
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
	result := db.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})

		return
	}

	db.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}*/
