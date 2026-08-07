package controller

import (
	"net/http"
	"plantao/internal/api/apierr"
	"plantao/internal/domain/role"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	service *role.RoleService
}

func NewRoleController(service *role.RoleService) *RoleController {
	return &RoleController{service: service}
}

func (c *RoleController) GetAll(ctx *gin.Context) {
	roles, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		apierr.Respond(ctx, err)
		return
	}

	type response struct {
		Id   string `json:"id"`
		Nome string `json:"nome"`
	}

	resp := make([]response, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, response{Id: r.Id.String(), Nome: r.Nome})
	}

	ctx.JSON(http.StatusOK, resp)
}
