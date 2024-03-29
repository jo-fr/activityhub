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
	if err := e.GetTX().Preload("Account").First(&source, "feed_url = ?", feedURL).Error; err != nil {
		return model.Feed{}, err
	}
	return source, nil
}

func (e *FeedRepository) GetFeedWithID(id string) (model.Feed, error) {
	var source model.Feed
	if err := e.GetTX().Preload("Account").First(&source, "id = ?", id).Error; err != nil {
		return model.Feed{}, err
	}
	return source, nil
}

func (e *FeedRepository) GetFeedWithAccountID(accountID string) (model.Feed, error) {
	var feed model.Feed
	if err := e.GetTX().Preload("Account").First(&feed, "account_id = ?", accountID).Error; err != nil {
		return model.Feed{}, err
	}
	return feed, nil
}

func (e *FeedRepository) FeedCount() (int64, error) {
	var count int64
	if err := e.GetTX().Model(&model.Feed{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *FeedRepository) StatusCount(accountID string) (int64, error) {
	var count int64
	if err := e.GetTX().Where("account_id = ?", accountID).Model(&model.Status{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *FeedRepository) ListFeeds(offset int, limit int) ([]model.Feed, error) {
	var feeds []model.Feed
	err := e.GetTX().
		Preload("Account").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&feeds).
		Error

	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (e *FeedRepository) ListStatusFromAccount(accountID string, offset int, limit int) ([]model.Status, error) {
	var status []model.Status
	err := e.GetTX().
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&status, "account_id = ?", accountID).
		Error

	if err != nil {
		return nil, err
	}
	return status, nil
}
