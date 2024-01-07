package activitypub

import (
	"context"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/internal/keys"
	"github.com/jo-fr/activityhub/backend/modules/activitypub/internal/repository"
	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
	"github.com/jo-fr/activityhub/backend/pkg/errutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// define errors
var (
	ErrActorNotFound = errutil.NewError(errutil.TypeNotFound, "actor not found")
)

func (h *Handler) GetActor(ctx context.Context, actor string) (acc models.Account, err error) {
	err = h.store.Execute(ctx, func(e *repository.ActivityHubRepository) error {
		acc, err = e.GetAccoutByUsername(actor)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrActorNotFound
			}
			return errors.Wrap(err, "failed to get actor from db")
		}

		return nil
	})

	return acc, err
}

func (h *Handler) CreateAccount(ctx context.Context, username string, name string, summary string) (acc models.Account, err error) {
	err = h.store.Execute(ctx, func(e *repository.ActivityHubRepository) error {
		keys, err := keys.GenerateRSAKeyPair(2048)
		if err != nil {
			return errors.Wrap(err, "failed to generate RSA key pair")
		}

		account := models.Account{
			PreferredUsername: username,
			Name:              name,
			Summary:           summary,
			PrivateKey:        []byte(keys.PrivKeyPEM),
			PublicKey:         []byte(keys.PubKeyPEM),
		}

		acc, err = e.CreateAccount(account)
		if err != nil {
			return errors.Wrap(err, "failed create account in db")
		}
		return nil
	})

	return acc, err
}
