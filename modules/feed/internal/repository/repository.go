package repository

import (
	"github.com/jo-fr/activityhub/pkg/store"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewFeedRepository),
	fx.Provide(store.NewStore[*FeedRepository]),
)

type FeedRepository struct {
	*store.Repository
}

func NewFeedRepository() *FeedRepository {
	return &FeedRepository{
		Repository: store.NewRepository(),
	}
}
