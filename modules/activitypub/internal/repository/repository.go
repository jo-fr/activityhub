package repository

import (
	"github.com/jo-fr/activityhub/pkg/store"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewActivityHubRepository),
	fx.Provide(store.NewStore[*ActivityHubRepository]),
)

type ActivityHubRepository struct {
	*store.Repository
}

func NewActivityHubRepository() *ActivityHubRepository {
	return &ActivityHubRepository{
		Repository: store.NewRepository(),
	}
}
