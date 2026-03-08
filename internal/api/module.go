package api

import (
	"net/http"
	"plantao/internal/api/controller"
	"plantao/internal/api/midware"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	fx.Provide(
		controller.NewPlantaoController,
		controller.NewColaboradorController,
		controller.NewUsuarioController,
		controller.NewAuthController,
		midware.NewAuthMidware,
		controller.NewModeloComunicacaoController,
		fx.Annotate(
			NewRouter,
			fx.As(new(http.Handler)),
		),
		NewServer,
	),
)
