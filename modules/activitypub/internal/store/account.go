package store

import (
	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (e *ActivityHubRepository) GetAccoutByUsername(username string) (models.Account, error) {
	var account models.Account

	err := e.GetTX().Where("preferred_username = ?", username).First(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (e *ActivityHubRepository) CreateAccount(account models.Account) (models.Account, error) {
	err := e.GetTX().Create(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (e *ActivityHubRepository) GetAccountByID(id string) (models.Account, error) {
	var account models.Account

	err := e.GetTX().Where("id = ?", id).First(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}
