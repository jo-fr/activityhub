package activitypub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub/internal/keys/httprequest"
	"github.com/jo-fr/activityhub/modules/activitypub/internal/store"
	"github.com/jo-fr/activityhub/modules/activitypub/models"

	"github.com/jo-fr/activityhub/pkg/externalmodel"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) ReceiveInboxActivity(ctx context.Context, activity externalmodel.Activity) error {
	return h.store.Execute(ctx, func(e *store.ActivityHubRepository) error {
		obj, ok := activity.Object.(string)
		if !ok {
			return errors.New("object is not a string")
		}

		accountName := getAccountFromURI(obj)
		account, err := e.GetAccoutByUsername(accountName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrActorNotFound
			}
			return errors.Wrap(err, "failed to get actor from db")
		}

		_, err = e.CreateFollow(account.ID, activity.Actor)
		if err != nil {
			return errors.Wrap(err, "failed to create follow")
		}

		if err := returnAcceptActivity(ctx, account, activity); err != nil {
			return errors.Wrap(err, "failed to return accept activity")
		}

		return nil
	})
}

func returnAcceptActivity(ctx context.Context, account models.Account, activity externalmodel.Activity) error {

	obj := activity.Object.(string)

	ma := externalmodel.Activity{
		Context: "https://www.w3.org/ns/activitystreams",
		ID:      activity.ID,
		Type:    "Accept",
		Actor:   obj,
		Object:  activity,
	}

	inboxURL, err := GetInboxURL(activity.Actor)
	if err != nil {
		return errors.Wrap(err, "failed to get inbox url")
	}

	json, err := json.Marshal(ma)
	if err != nil {
		return errors.Wrap(err, "failed to marshal activity")
	}

	req, err := httprequest.New(http.MethodPost, inboxURL, bytes.NewBuffer(json))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	if err := req.Sign(account.PrivateKey, obj); err != nil {
		return errors.Wrap(err, "failed to sign request")
	}

	resp, err := req.Do()
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	if !httputil.StatusOK(resp.StatusCode) {
		errBody, err := httputil.UnmarshalBody[map[string]any](resp.Body)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal error response body. Statuscode: %v", resp.StatusCode)
		}
		return fmt.Errorf("received status code %d. Response: %s", resp.StatusCode, errBody)
	}

	return nil
}

func GetInboxURL(actorURI string) (string, error) {
	req, err := httprequest.New(http.MethodGet, actorURI, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}

	resp, err := req.Do()
	if err != nil {
		return "", errors.Wrap(err, "failed to get actor")
	}

	responseMap, err := httputil.UnmarshalBody[map[string]any](resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response body")
	}

	inbox, ok := responseMap["inbox"].(string)
	if !ok {
		return "", errors.New("failed to get inbox from actor")
	}

	return inbox, nil
}

// getAccountFromURI returns the account name from an URI
// e.g. "https://example.com/users/account" -> "account"
func getAccountFromURI(uri string) string {
	// Split the URI by "/"
	parts := strings.Split(uri, "/")
	// Return the last part of the URI
	return parts[len(parts)-1]
}
