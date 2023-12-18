package repository

import (
	"github.com/jo-fr/activityhub/pkg/store"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(store.NewStore[FeedRepository]),
)

type FeedRepository struct {
	*store.Repository
}

func (e FeedRepository) WithRepository() FeedRepository {
	e.Repository = store.NewRepository()
	return e
}
