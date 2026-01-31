package api

import (
	"net/http"
	"plantao/internal/api/controller"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	fx.Provide(
		controller.NewPlantaoController,
		fx.Annotate(
			NewRouter,
			fx.As(new(http.Handler)),
		),
		NewServer,
	),
)
