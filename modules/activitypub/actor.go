package activitypub

import (
	"context"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/jo-fr/activityhub/pkg/errutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// define errors
var (
	ErrActorNotFound = errutil.NewError(errutil.TypeNotFound, "actor not found")
)

func (h *Handler) GetActor(ctx context.Context, actor string) (models.Account, error) {
	account, err := h.store.GetAccoutByUsername(ctx, actor)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Account{}, ErrActorNotFound
		}
		return models.Account{}, errors.Wrap(err, "failed to get actor from db")
	}

	return account, nil
}
