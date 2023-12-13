package store

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

func (e *FeedRepository) ListSourceFeeds() ([]model.SourceFeed, error) {
	var sources []model.SourceFeed
	if err := e.GetTX().Order("created_at DESC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}
