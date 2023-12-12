package store

import (
	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (e *Executer) GetAccoutByUsername(username string) (models.Account, error) {
	var account models.Account

	err := e.tx.Where("preferred_username = ?", username).First(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (e *Executer) CreateAccount(account models.Account) (models.Account, error) {
	err := e.tx.Create(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (e *Executer) GetAccountByID(id string) (models.Account, error) {
	var account models.Account

	err := e.tx.Where("id = ?", id).First(&account).Error
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}
