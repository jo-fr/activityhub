package database

import (
	"context"

	"gorm.io/gorm"
)

// Store is a generic data store that works with types implementing the IRepository interface.
type Store[e IRepository] struct {
	db *Database
}

// NewStore creates a new instance of the Store with the provided database connection.
func NewStore[T IRepository](db *Database) *Store[T] {
	return &Store[T]{
		db: db,
	}
}

// IRepository is an interface defining methods for managing transactions with a database.
type IRepository interface {
	GetTX() *gorm.DB
	SetTX(tx *gorm.DB)
}

// Execute executes a function within a database transaction. It either uses the transaction
// provided in the context or creates a new one. It rolls back the transaction if an error occurs.
func (s *Store[T]) Execute(ctx context.Context, f func(e T) error) error {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		tx = s.db.DB.WithContext(ctx).Begin()
		defer tx.Rollback()
	}

	var repository T
	repository.SetTX(tx)

	if err := f(repository); err != nil {
		return err
	}

	// Only commit if tx is not passed in context
	if !ok {
		return tx.Commit().Error
	}

	return nil
}

// Repository is a concrete implementation of the IRepository interface.
type Repository struct {
	tx *gorm.DB
}

// NewRepository creates a new instance of the Repository.
func NewRepository() *Repository {
	return &Repository{}
}

// GetTX returns the current database transaction associated with the repository.
func (e *Repository) GetTX() *gorm.DB {
	return e.tx
}

// SetTX sets the database transaction for the repository.
func (e *Repository) SetTX(tx *gorm.DB) {
	e.tx = tx
}
