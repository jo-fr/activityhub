package repository

import (
	"github.com/jo-fr/activityhub/modules/feed/model"
)

func (e *FeedRepository) GetLatestStatusFromSourceFeed(accountID string) (model.Status, error) {
	var status model.Status
	err := e.GetTX().
		Order("created_at DESC").
		Take(&status, "account_id = ?", accountID).
		Error

	if err != nil {
		return model.Status{}, err
	}
	return status, nil
}

func (e *FeedRepository) CreateStatus(status model.Status) (model.Status, error) {
	if err := e.GetTX().Create(&status).Error; err != nil {
		return model.Status{}, err
	}
	return status, nil
}
