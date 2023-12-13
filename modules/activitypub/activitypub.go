package activitypub

import (
	"github.com/jo-fr/activityhub/modules/activitypub/internal/store"
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/jo-fr/activityhub/pkg/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	store.Module,
	fx.Provide(ProvideHandler),
)

type Handler struct {
	hostURL string
	store   *database.Store[*store.ActivityHubRepository]
}

func ProvideHandler(config config.Config, store *database.Store[*store.ActivityHubRepository]) *Handler {

	return &Handler{
		hostURL: config.HostURL,
		store:   store,
	}
}
