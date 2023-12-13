package database

import (
	"context"

	"gorm.io/gorm"
)

type Store[e IRepository] struct {
	db *Database
}

func NewStore[T IRepository](db *Database) *Store[T] {
	return &Store[T]{
		db: db,
	}
}

type IRepository interface {
	GetTX() *gorm.DB
	SetTX(tx *gorm.DB)
}

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

	// only commit if tx is not passed in context
	if !ok {
		return tx.Commit().Error
	}

	return nil

}

type Repository struct {
	tx *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (e *Repository) GetTX() *gorm.DB {
	return e.tx

}

func (e *Repository) SetTX(tx *gorm.DB) {
	e.tx = tx
}
