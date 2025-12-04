package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main.go/types"
)

type UserController struct {
	repository types.UserRepository
}

func NewUserController(repository types.UserRepository) *UserController {
	return &UserController{
		repository: repository,
	}
}

func (c *UserController) RegisterRoutes(r *gin.Engine) {
	r.POST("/users/:id", c.userRegister)
	r.GET("/users", c.userGet)
	r.GET("/users/:id", c.userGetById)
	r.PUT("/users/:id", c.userUpdate)
	r.DELETE("/users/:id", c.userDelete)
}

func (c *UserController) userRegister(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	var user types.UserRequest

	g.Bind(&user)

	userSaved, err := c.repository.CreateUser(idConverted, user)

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

func (c *UserController) userGet(g *gin.Context) {
	var user types.UserRequest

	g.Bind(&user)

	users, err := c.repository.GetUsers()

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	g.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (c *UserController) userGetById(g *gin.Context) {
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

func (c *UserController) userUpdate(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	var user types.UserRequest
	g.Bind(&user)

	err := c.repository.UpdateUser(idConverted, user)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Usuário atualizado com sucesso",
	})
}

func (c *UserController) userDelete(g *gin.Context) {
	idConverted, errId := uuid.Parse(g.Param("id"))

	if errId != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": errId.Error(),
		})
		return
	}

	err := c.repository.DeleteUserById(idConverted)

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
