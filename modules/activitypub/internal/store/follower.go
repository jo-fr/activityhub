package store

import (
	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (e *Executer) CreateFollow(accountIDFollowed string, accountURIFollowing string) (models.Follower, error) {
	follower := models.Follower{
		AccountIDFollowed:   accountIDFollowed,
		AccountURIFollowing: accountURIFollowing,
	}

	if err := e.tx.Create(&follower).Error; err != nil {
		return models.Follower{}, err
	}

	return follower, nil
}

func (e *Executer) GetFollowersOfAccount(accountID string) ([]models.Follower, error) {
	var followers []models.Follower

	err := e.tx.Where("account_id_followed = ?", accountID).Find(&followers).Error
	if err != nil {
		return nil, err
	}

	return followers, nil
}
