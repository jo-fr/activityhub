package activitypub

import (
	"github.com/jo-fr/activityhub/backend/modules/activitypub/internal/repository"
	"github.com/jo-fr/activityhub/backend/pkg/config"
	"github.com/jo-fr/activityhub/backend/pkg/store"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repository.Module,
	fx.Provide(ProvideHandler),
	fx.Provide(NewConsumer),
	fx.Invoke(Subscribe),
)

type Handler struct {
	host  string
	store *store.Store[repository.ActivityHubRepository]
}

func ProvideHandler(config config.Config, store *store.Store[repository.ActivityHubRepository]) *Handler {

	return &Handler{
		host:  config.Host,
		store: store,
	}
}
