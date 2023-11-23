package store

import (
	"context"

	"github.com/jo-fr/activityhub/modules/feed/model"
)

func (s *Store) CreateSourceFeed(ctx context.Context, source model.SourceFeed) (model.SourceFeed, error) {
	if err := s.db.WithContext(ctx).Create(&source).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (s *Store) GetSourceFeedWithFeedURL(ctx context.Context, feedURL string) (model.SourceFeed, error) {
	var source model.SourceFeed
	if err := s.db.WithContext(ctx).Where("url = ?", feedURL).First(&source).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}
