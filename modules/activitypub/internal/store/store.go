package store

import (
	"context"

	"github.com/jo-fr/activityhub/pkg/database"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(NewStore),
)

func NewStore(db *database.Database) *Store {
	return &Store{
		db: db,
	}
}

type Executer struct {
	tx *gorm.DB
}

type Store struct {
	db *database.Database
}

func (s *Store) Execute(ctx context.Context, f func(e *Executer) error) error {
	tx := s.db.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := f(&Executer{tx: tx}); err != nil {
		return err
	}

	return tx.Commit().Error

}
