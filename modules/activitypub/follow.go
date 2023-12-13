package activitypub

import (
	"context"

	"github.com/jo-fr/activityhub/modules/activitypub/internal/store"
	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) GetFollowers(ctx context.Context, actorname string) (follower []models.Follower, err error) {

	err = h.store.Execute(ctx, func(e *store.ActivityHubRepository) error {
		account, err := e.GetAccoutByUsername(actorname)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrActorNotFound
			}
			return errors.Wrap(err, "failed to get actor from db")
		}

		follower, err = e.GetFollowersOfAccount(account.ID)
		if err != nil {
			return err
		}
		return nil
	})

	return follower, err
}
