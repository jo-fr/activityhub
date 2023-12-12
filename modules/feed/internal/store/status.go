package store

import (
	"github.com/jo-fr/activityhub/modules/feed/model"
)

func (e *Executer) GetLatestStatusFromSourceFeed(accountID string) (model.Status, error) {
	var status model.Status
	err := e.tx.
		Order("created_at DESC").
		Take(&status, "account_id = ?", accountID).
		Error

	if err != nil {
		return model.Status{}, err
	}
	return status, nil
}

func (e *Executer) CreateStatus(status model.Status) (model.Status, error) {
	if err := e.tx.Create(&status).Error; err != nil {
		return model.Status{}, err
	}
	return status, nil
}
