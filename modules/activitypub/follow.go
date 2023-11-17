package activitypub

import (
	"context"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) GetFollowers(ctx context.Context, actorname string) ([]models.Follower, error) {

	account, err := h.store.GetAccoutByUsername(ctx, actorname)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActorNotFound
		}
		return nil, errors.Wrap(err, "failed to get actor from db")
	}

	followers, err := h.store.GetFollowersOfAccount(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}
