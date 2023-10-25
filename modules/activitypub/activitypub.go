package activitypub

import (
	"github.com/jo-fr/activityhub/pkg/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(ProvideHandler),
)

type Handler struct {
	hostURL string
}

func ProvideHandler(config config.Config) *Handler {
	return &Handler{
		hostURL: config.HostURL,
	}
}
