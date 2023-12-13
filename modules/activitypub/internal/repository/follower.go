package repository

import (
	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (e *ActivityHubRepository) CreateFollow(accountIDFollowed string, accountURIFollowing string) (models.Follower, error) {
	follower := models.Follower{
		AccountIDFollowed:   accountIDFollowed,
		AccountURIFollowing: accountURIFollowing,
	}

	if err := e.GetTX().Create(&follower).Error; err != nil {
		return models.Follower{}, err
	}

	return follower, nil
}

func (e *ActivityHubRepository) GetFollowersOfAccount(accountID string) ([]models.Follower, error) {
	var followers []models.Follower

	err := e.GetTX().Where("account_id_followed = ?", accountID).Find(&followers).Error
	if err != nil {
		return nil, err
	}

	return followers, nil
}
