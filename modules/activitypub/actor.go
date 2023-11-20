package activitypub

import (
	"context"
	"fmt"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub/internal/keys"
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

func (h *Handler) CreateAccount(ctx context.Context, username string) (models.Account, error) {

	keys, err := keys.GenerateRSAKeyPair(2048)
	if err != nil {
		return models.Account{}, errors.Wrap(err, "failed to generate RSA key pair")
	}

	account := models.Account{
		PreferredUsername: username,
		Name:              strings.Title(username),
		Summary:           fmt.Sprintf("This is the mastodon Account of %s", username),
		PrivateKey:        []byte(keys.PrivKeyPEM),
		PublicKey:         []byte(keys.PubKeyPEM),
	}

	account, err = h.store.CreateAccount(ctx, account)
	if err != nil {
		return models.Account{}, errors.Wrap(err, "failed create account in db")
	}
	return account, nil
}
