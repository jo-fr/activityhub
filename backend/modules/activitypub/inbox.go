package activitypub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/internal/keys/httprequest"
	"github.com/jo-fr/activityhub/backend/modules/activitypub/internal/repository"
	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"

	"github.com/jo-fr/activityhub/backend/pkg/errutil"
	"github.com/jo-fr/activityhub/backend/pkg/externalmodel"
	"github.com/jo-fr/activityhub/backend/pkg/util/httputil"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrUnsupportedActivityType = errutil.NewError(errutil.TypeBadRequest, "unsupported activity type")
)

func (h *Handler) ReceiveInboxActivity(ctx context.Context, activity externalmodel.Activity) error {
	var account models.Account
	var object string
	err := h.store.Execute(ctx, func(e *repository.ActivityHubRepository) error {
		actor, _object, err := extractActorAndObject(ctx, activity)
		if err != nil {
			return errors.Wrap(err, "failed to extract actor and object")
		}
		object = _object
		accountName := getAccountFromURI(object)
		account, err = e.GetAccoutByUsername(accountName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrActorNotFound
			}
			return errors.Wrap(err, "failed to get actor from db")
		}

		switch activity.Type {
		case "Follow":
			if _, err := e.CreateFollow(account.ID, actor); err != nil {
				return errors.Wrap(err, "failed to create follow")
			}
		case "Undo":
			if err := e.DeleteFollow(account.ID, actor); err != nil {
				return errors.Wrap(err, "failed to delete follow")
			}
		default:
			return ErrUnsupportedActivityType
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err = returnAcceptActivity(ctx, account, object, activity); err != nil {
		return errors.Wrap(err, "failed to return accept activity")
	}

	return nil
}

func extractActorAndObject(ctx context.Context, activity externalmodel.Activity) (actor string, object string, err error) {

	switch activity.Type {
	case "Follow":
		actor = activity.Actor
		object, ok := activity.Object.(string)
		if !ok {
			return "", "", errors.New("object must be string for follow activity")
		}

		return actor, object, nil

	case "Undo":
		actor = activity.Actor
		var nestedActivity externalmodel.Activity
		if err := mapstructure.Decode(activity.Object, &nestedActivity); err != nil {
			return "", "", errors.Wrap(err, "failed to decode nested activity")
		}

		if nestedActivity.Type != "Follow" {
			return "", "", ErrUnsupportedActivityType
		}

		obj, ok := nestedActivity.Object.(string)
		if !ok {
			return "", "", errors.New("object must be string for follow activity")
		}

		return actor, obj, nil
	default:
		return "", "", ErrUnsupportedActivityType
	}
}

func returnAcceptActivity(ctx context.Context, account models.Account, actor string, activity externalmodel.Activity) error {

	ma := externalmodel.Activity{
		Context: "https://www.w3.org/ns/activitystreams",
		ID:      activity.ID,
		Type:    "Accept",
		Actor:   actor,
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

	if err := req.Sign(account.PrivateKey, actor); err != nil {
		return errors.Wrap(err, "failed to sign request")
	}

	resp, err := req.Do()
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if !httputil.StatusOK(resp.StatusCode) {
		errBody, err := httputil.UnmarshalResponsetBody[map[string]any](resp)
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

	responseMap, err := httputil.UnmarshalResponsetBody[map[string]any](resp)
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
