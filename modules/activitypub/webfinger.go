package activitypub

import (
	"fmt"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/jo-fr/activityhub/pkg/errutil"
)

var (
	ErrWrongURI    = errutil.NewError(errutil.TypeInvalidRequestBody, "uri of actor and host do not match")
	ErrWrongFormat = errutil.NewError(errutil.TypeBadRequest, "no acct: in resource")
)

func (h *Handler) GetWebfinger(resource string) (models.Webfinger, error) {

	store := h.db

	username, err := validateAndExtractActor(h.hostURL, resource)
	if err != nil {
		return models.Webfinger{}, err
	}

	var acc models.Account
	err = store.First(&acc, "preferred_username = ?", username).Error
	if err != nil {
		return models.Webfinger{}, err
	}

	webfinger := models.Webfinger{
		Subject: resource,
		Links: []models.Links{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("%s/%s", h.hostURL, username),
			},
		},
	}

	return webfinger, nil
}

func validateAndExtractActor(hostURL string, resource string) (string, error) {
	if !strings.Contains(resource, "acct:") {
		return "", ErrWrongFormat
	}

	actor := strings.Split(resource, "acct:")[1]
	username := strings.Split(actor, "@")[0]
	uri := strings.Split(actor, "@")[1]

	if uri != hostURL {
		return "", ErrWrongURI
	}

	return username, nil
}
