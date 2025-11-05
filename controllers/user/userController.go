package user

import (
	"acommerce-api/types/dtos"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/hello", h.hello)
	r.POST("/user", h.createUser)
}

func (h *Handler) hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, User!"})
}

func (h *Handler) createUser(c *gin.Context) {
	var userDTO dtos.UserDTO

	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}
