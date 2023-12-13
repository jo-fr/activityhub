package store

import (
	"github.com/jo-fr/activityhub/pkg/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewActivityHubRepository),
	fx.Provide(database.NewStore[*ActivityHubRepository]),
)

type ActivityHubRepository struct {
	*database.Repository
}

func NewActivityHubRepository() *ActivityHubRepository {
	return &ActivityHubRepository{
		Repository: database.NewRepository(),
	}
}
