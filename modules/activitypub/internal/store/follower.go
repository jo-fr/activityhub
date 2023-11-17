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
