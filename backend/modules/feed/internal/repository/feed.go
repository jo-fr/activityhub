package repository

import (
	"github.com/jo-fr/activityhub/backend/modules/feed/model"
)

func (e *FeedRepository) CreateFeed(source model.Feed) (model.Feed, error) {
	if err := e.GetTX().Create(&source).Error; err != nil {
		return model.Feed{}, err
	}
	return source, nil
}

func (e *FeedRepository) GetFeedWithFeedURL(feedURL string) (model.Feed, error) {
	var source model.Feed
	if err := e.GetTX().First(&source, "feed_url = ?", feedURL).Error; err != nil {
		return model.Feed{}, err
	}
	return source, nil
}

func (e *FeedRepository) GetFeedWithID(id string) (model.Feed, error) {
	var source model.Feed
	if err := e.GetTX().First(&source, "id = ?", id).Error; err != nil {
		return model.Feed{}, err
	}
	return source, nil
}

func (e *FeedRepository) CountFeeds() (int64, error) {
	var count int64
	if err := e.GetTX().Model(&model.Feed{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *FeedRepository) ListFeeds(offset int, limit int) ([]model.Feed, error) {
	var sources []model.Feed
	err := e.GetTX().
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sources).
		Error

	if err != nil {
		return nil, err
	}
	return sources, nil
}
