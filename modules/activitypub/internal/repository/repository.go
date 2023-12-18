package repository

import (
	"github.com/jo-fr/activityhub/pkg/store"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(store.NewStore[ActivityHubRepository]),
)

type ActivityHubRepository struct {
	*store.Repository
}

func (e ActivityHubRepository) WithRepository() ActivityHubRepository {
	e.Repository = store.NewRepository()
	return e
}
