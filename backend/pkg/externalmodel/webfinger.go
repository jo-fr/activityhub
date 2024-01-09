package externalmodel

import (
	"fmt"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
)

type Webfinger struct {
	Subject string   `json:"subject,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
	Links   []Links  `json:"links,omitempty"`
}
type Links struct {
	Rel      string `json:"rel,omitempty"`
	Type     string `json:"type,omitempty"`
	Href     string `json:"href,omitempty"`
	Template string `json:"template,omitempty"`
}

func ExternalWebfinger(hostURL string, resource string, acc models.Account) Webfinger {
	return Webfinger{
		Subject: resource,
		Links: []Links{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s/ap/%s", hostURL, acc.PreferredUsername),
			},
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: fmt.Sprintf("https://%s/ap/%s/tese", hostURL, acc.PreferredUsername),
			},
		},
	}
}
