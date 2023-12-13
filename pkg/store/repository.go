package store

import (
	"context"

	"gorm.io/gorm"
)

// Repository is a concrete implementation of the IRepository interface.
type Repository struct {
	tx *gorm.DB
}

// IRepository is an interface defining methods for managing transactions with a database.
type IRepository interface {
	GetTX() *gorm.DB
	SetTX(tx *gorm.DB)
	GetCtxWithTx(ctx context.Context) context.Context
}

// NewRepository creates a new instance of the Repository.
func NewRepository() *Repository {
	return &Repository{}
}

// GetCtxWithTx returns a new context with the current database transaction associated with the repository.
func (e *Repository) GetCtxWithTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextWithTxKey{}, e.GetTX())
}

// GetTX returns the current database transaction associated with the repository.
func (e *Repository) GetTX() *gorm.DB {
	if e == nil {
		return nil
	}

	return e.tx
}

// SetTX sets the database transaction for the repository.
func (e *Repository) SetTX(tx *gorm.DB) {
	if e == nil {
		e = &Repository{}
	}
	e.tx = tx
}
