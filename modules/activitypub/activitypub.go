package activitypub

import (
	"github.com/jo-fr/activityhub/modules/activitypub/internal/store"
	"github.com/jo-fr/activityhub/pkg/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	store.Module,
	fx.Provide(ProvideHandler),
)

type Handler struct {
	hostURL string
	store   *store.Store
}

func ProvideHandler(config config.Config, store *store.Store) *Handler {

	return &Handler{
		hostURL: config.HostURL,
		store:   store,
	}
}
