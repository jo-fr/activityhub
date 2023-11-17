package store

import (
	"context"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (s *Store) CreateFollow(ctx context.Context, accountIDFollowed string, accountURIFollowing string) (models.Follower, error) {
	follower := models.Follower{
		AccountIDFollowed:   accountIDFollowed,
		AccountURIFollowing: accountURIFollowing,
	}

	if err := s.db.WithContext(ctx).Create(&follower).Error; err != nil {
		return models.Follower{}, err
	}

	return follower, nil
}

func (s *Store) GetFollowersOfAccount(ctx context.Context, accountID string) ([]models.Follower, error) {
	var followers []models.Follower

	err := s.db.WithContext(ctx).Where("account_id_followed = ?", accountID).Find(&followers).Error
	if err != nil {
		return nil, err
	}

	return followers, nil
}
