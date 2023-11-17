package store

import (
	"github.com/jo-fr/activityhub/pkg/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewStore),
)

func NewStore(db *database.Database) *Store {
	return &Store{
		db: db,
	}
}

type Store struct {
	db *database.Database
}
