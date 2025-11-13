package controllers

import (
	"github.com/gin-gonic/gin"
	"main.go/types"
)

type ColaboradorController struct {
	repository types.ColaboradorRepository
}

func NewColaboradorController(repository types.ColaboradorRepository) *ColaboradorController {
	return &ColaboradorController{
		repository: repository,
	}
}

func (c *ColaboradorController) RegisterRoutes(r *gin.Engine) {

}
