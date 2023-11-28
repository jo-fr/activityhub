package store

import (
	"context"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (s *Store) GetAccoutByUsername(ctx context.Context, username string) (models.Account, error) {
	var account models.Account

	err := s.db.WithContext(ctx).Where("preferred_username = ?", username).First(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (s *Store) CreateAccount(ctx context.Context, account models.Account) (models.Account, error) {
	err := s.db.WithContext(ctx).Create(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (s *Store) GetAccountByID(ctx context.Context, id string) (models.Account, error) {
	var account models.Account

	err := s.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}
