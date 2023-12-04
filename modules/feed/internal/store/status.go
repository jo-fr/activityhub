package store

import (
	"context"

	"github.com/jo-fr/activityhub/modules/feed/model"
)

func (s *Store) GetLatestStatusFromSourceFeed(ctx context.Context, accountID string) (model.Status, error) {
	var status model.Status
	err := s.db.
		WithContext(ctx).
		Order("created_at DESC").
		Take(&status, "account_id = ?", accountID).
		Error

	if err != nil {
		return model.Status{}, err
	}
	return status, nil
}

func (s *Store) CreateStatus(ctx context.Context, status model.Status) (model.Status, error) {
	if err := s.db.WithContext(ctx).Create(&status).Error; err != nil {
		return model.Status{}, err
	}
	return status, nil
}
