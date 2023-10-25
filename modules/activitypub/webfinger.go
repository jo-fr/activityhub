package activitypub

import (
	"fmt"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub/internal/models"
)

func (h *Handler) GetWebfinger(resource string) (models.Webfinger, error) {

	username, err := extractActor(resource)
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

func extractActor(resource string) (string, error) {
	if !strings.Contains(resource, "acct:") {
		return "", fmt.Errorf("no acct: in resource")
	}

	actorEmail := strings.Split(resource, "acct:")[1]
	actor := strings.Split(actorEmail, "@")[0]

	return actor, nil
}
