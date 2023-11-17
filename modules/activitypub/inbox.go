package activitypub

import (
	"context"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) ReceiveInboxActivity(ctx context.Context, actor string, object string, activityType string) (models.Follower, error) {

	accountName := getAccountFromURI(object)
	account, err := h.store.GetAccoutByUsername(ctx, accountName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Follower{}, ErrActorNotFound
		}
		return models.Follower{}, errors.Wrap(err, "failed to get actor from db")
	}

	follow, err := h.store.CreateFollow(ctx, account.ID, object)
	if err != nil {
		return models.Follower{}, errors.Wrap(err, "failed to create follow")
	}

	return follow, nil
}

// getAccountFromURI returns the account name from an URI
// e.g. "https://example.com/users/account" -> "account"
func getAccountFromURI(uri string) string {
	// Split the URI by "/"
	parts := strings.Split(uri, "/")
	// Return the last part of the URI
	return parts[len(parts)-1]
}
