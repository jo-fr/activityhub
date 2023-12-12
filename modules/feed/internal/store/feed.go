package store

import (
	"github.com/jo-fr/activityhub/modules/feed/model"
)

func (e *Executer) CreateSourceFeed(source model.SourceFeed) (model.SourceFeed, error) {
	if err := e.tx.Create(&source).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (e *Executer) GetSourceFeedWithFeedURL(feedURL string) (model.SourceFeed, error) {
	var source model.SourceFeed
	if err := e.tx.First(&source, "feed_url = ?", feedURL).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (e *Executer) GetSourceFeedWithID(id string) (model.SourceFeed, error) {
	var source model.SourceFeed
	if err := e.tx.First(&source, "id = ?", id).Error; err != nil {
		return model.SourceFeed{}, err
	}
	return source, nil
}

func (e *Executer) ListSourceFeeds() ([]model.SourceFeed, error) {
	var sources []model.SourceFeed
	if err := e.tx.Order("created_at DESC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}
