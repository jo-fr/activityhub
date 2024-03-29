package store

import (
	"context"

	"github.com/jo-fr/activityhub/backend/pkg/database"
	"gorm.io/gorm"
)

// contextWithTxKey is the key used to store the database transaction in the context.
type contextWithTxKey struct{}

// Store is a generic data store that works with types implementing the IRepository interface.
type Store[T IRepository[T]] struct {
	db *database.Database
}

// NewStore creates a new instance of the Store with the provided database connection.
func NewStore[T IRepository[T]](db *database.Database) *Store[T] {
	return &Store[T]{
		db: db,
	}
}

// Execute executes a repository functions within one database transaction. It either uses the transaction
// provided in the context or creates a new one. It rolls back the transaction if an error occurs.
func (s *Store[T]) Execute(ctx context.Context, f func(e *T) error) error {
	tx, ok := ctx.Value(contextWithTxKey{}).(*gorm.DB)
	if !ok {
		tx = s.db.DB.WithContext(ctx).Begin()
		defer tx.Rollback()
	}

	var rep T
	rep = rep.WithRepository()
	rep.SetTX(tx)

	if err := f(&rep); err != nil {
		return err
	}

	// Only commit if tx is not passed in context
	if !ok {
		return tx.Commit().Error
	}

	return nil
}
