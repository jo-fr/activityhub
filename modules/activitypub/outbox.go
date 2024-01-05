package activitypub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jo-fr/activityhub/modules/activitypub/internal/keys/httprequest"
	"github.com/jo-fr/activityhub/modules/activitypub/internal/repository"
	"github.com/jo-fr/activityhub/pkg/externalmodel"
	"github.com/jo-fr/activityhub/pkg/util"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
	"github.com/pkg/errors"
)

func (h *Handler) SendPostToFollowers(ctx context.Context, sendingActorID string, content string) error {
	return h.store.Execute(ctx, func(e *repository.ActivityHubRepository) error {
		followers, err := e.GetFollowersOfAccount(sendingActorID)
		if err != nil {
			return errors.Wrap(err, "failed to get followers of account")
		}

		for _, follower := range followers {
			if err := h.SendPost(ctx, sendingActorID, follower.AccountURIFollowing, content); err != nil {
				return errors.Wrap(err, "failed to send post")
			}
		}
		return nil
	})
}

func (h *Handler) SendPost(ctx context.Context, sendingActorID string, sendToURI string, content string) error {
	return h.store.Execute(ctx, func(e *repository.ActivityHubRepository) error {
		account, err := e.GetAccountByID(sendingActorID)
		if err != nil {
			return errors.Wrap(err, "failed to get account by id")
		}

		inboxURL, err := GetInboxURL(sendToURI)
		if err != nil {
			return errors.Wrap(err, "failed to get inbox url")
		}

		u := uuid.NewString()

		actorURI := h.builtAccountURI(account.PreferredUsername)
		activity := externalmodel.Activity{
			Context:   "https://www.w3.org/ns/activitystreams",
			ID:        fmt.Sprintf("%s#%s", actorURI, u),
			Type:      "Create",
			Actor:     actorURI,
			Published: time.Now().UTC().Format(time.RFC3339),
			Object: externalmodel.Activity{
				ID:        fmt.Sprintf("%s#%s", actorURI, u),
				Type:      "Note",
				Content:   content,
				Published: time.Now().UTC().Format(time.RFC3339),
				Sensitive: util.ToPointer(false),
			},
			To: []string{"https://www.w3.org/ns/activitystreams#Public"},
		}

		json, err := json.Marshal(activity)
		if err != nil {
			return errors.Wrap(err, "failed to marshal activity")
		}

		req, err := httprequest.New(http.MethodPost, inboxURL, bytes.NewBuffer(json))
		if err != nil {
			return errors.Wrap(err, "failed to create request")
		}

		if err := req.Sign(account.PrivateKey, actorURI); err != nil {
			return errors.Wrap(err, "failed to sign request")
		}

		resp, err := req.Do()
		if err != nil {
			return errors.Wrap(err, "failed to do request")
		}

		if !httputil.StatusOK(resp.StatusCode) {
			return errors.Errorf("request failed with status %s", resp.Status)
		}

		return nil
	})

}

func (h *Handler) builtAccountURI(username string) string {
	return fmt.Sprintf("https://%s/users/%s", h.hostURL, username)
}
