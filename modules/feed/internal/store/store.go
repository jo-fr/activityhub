package store

import (
	"github.com/jo-fr/activityhub/pkg/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(database.NewStore[*FeedRepository]),
)

type FeedRepository struct {
	*database.Repository
}
