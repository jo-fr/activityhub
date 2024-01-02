package repository

import (
	"github.com/jo-fr/activityhub/modules/feed/model"
)

func (e *FeedRepository) CreateSourceFeed(source model.SourceFeed) (model.SourceFeed, error) {
	if err := e.GetTX().Create(&source).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (e *FeedRepository) GetSourceFeedWithFeedURL(feedURL string) (model.SourceFeed, error) {
	var source model.SourceFeed
	if err := e.GetTX().First(&source, "feed_url = ?", feedURL).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (e *FeedRepository) GetSourceFeedWithID(id string) (model.SourceFeed, error) {
	var source model.SourceFeed
	if err := e.GetTX().First(&source, "id = ?", id).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (e *FeedRepository) CountSourceFeeds() (int64, error) {
	var count int64
	if err := e.GetTX().Model(&model.SourceFeed{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *FeedRepository) ListSourceFeeds(offset int, limit int) ([]model.SourceFeed, error) {
	var sources []model.SourceFeed
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
