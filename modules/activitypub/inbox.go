package activitypub

import (
	"context"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (h *Handler) ReceiveInboxActivity(ctx context.Context, actor string, object string, activityType string) (models.Follower, error) {

	accountName := getAccountFromURI(object)

	var account models.Account
	if err := h.db.First(&account, "preferred_username = ?", accountName).Error; err != nil {
		return models.Follower{}, err
	}

	follower := models.Follower{
		AccountIDFollowed:   account.ID,
		AccountURIFollowing: actor,
	}

	if err := h.db.Create(&follower).Error; err != nil {
		return models.Follower{}, err
	}

	return follower, nil
}
func getAccountFromURI(uri string) string {
	// Split the URI by "/"
	parts := strings.Split(uri, "/")
	// Return the last part of the URI
	return parts[len(parts)-1]
}
