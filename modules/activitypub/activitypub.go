package activitypub

import (
	"github.com/jo-fr/activityhub/modules/activitypub/internal/repository"
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/jo-fr/activityhub/pkg/store"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repository.Module,
	fx.Provide(ProvideHandler),
)

type Handler struct {
	hostURL string
	store   *store.Store[repository.ActivityHubRepository]
}

func ProvideHandler(config config.Config, store *store.Store[repository.ActivityHubRepository]) *Handler {

	return &Handler{
		hostURL: config.HostURL,
		store:   store,
	}
}
