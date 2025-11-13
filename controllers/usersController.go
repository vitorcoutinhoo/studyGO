package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main.go/types"
)

type Controller struct {
	repository types.UserRepository
}

func NewUserController(repository types.UserRepository) *Controller {
	return &Controller{
		repository: repository,
	}
}

func (c *Controller) RegisterRoutes(r *gin.Engine) {
	r.POST("/users", c.userRegister)
	r.GET("/users", c.userGet)
	r.GET("/users/:id", c.userGetById)
	r.PUT("/users/:id", c.userUpdate)
	r.DELETE("/users/:id", c.userDelete)
}

func (c *Controller) userRegister(g *gin.Context) {
	var user types.UserRequest

	g.Bind(&user)

	userSaved, err := c.repository.CreateUser(user)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	g.JSON(http.StatusOK, gin.H{
		"user": userSaved,
	})
}

func (c *Controller) userGet(g *gin.Context) {
	var user types.UserRequest

	g.Bind(&user)

	userCreated, err := c.repository.GetUsers()

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	g.JSON(http.StatusOK, gin.H{
		"users": userCreated,
	})
}

func (c *Controller) userGetById(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	user, err := c.repository.GetUserById(idConverted)

	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (c *Controller) userUpdate(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": "Sexo" + errId.Error(),
		})
		return
	}

	var user types.UserRequest
	g.Bind(&user)

	err := c.repository.UpdateUser(idConverted, user)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": "Sexo2" + err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Usuário atualizado com sucesso",
	})
}

func (c *Controller) userDelete(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	err := c.repository.DeletUserById(idConverted)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Usuário deletado com sucesso",
	})
}
