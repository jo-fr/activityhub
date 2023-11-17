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
