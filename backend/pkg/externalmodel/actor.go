package externalmodel

import (
	"fmt"
	"time"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
)

type Actor struct {
	Context           []string  `json:"@context"`
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Following         string    `json:"following"`
	Followers         string    `json:"followers"`
	PreferredUsername string    `json:"preferredUsername"`
	Name              string    `json:"name"`
	Summary           string    `json:"summary"`
	URL               string    `json:"url"`
	Published         string    `json:"published"`
	Inbox             string    `json:"inbox"`
	PublicKey         PublicKey `json:"publicKey"`
	Attachment        []string  `json:"attachment"`
}
type PublicKey struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

func ExternalActor(hostURL string, acc models.Account) Actor {

	username := acc.PreferredUsername

	return Actor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		ID:                fmt.Sprintf("https://%s/users/%s", hostURL, username),
		Type:              "Service",
		Following:         fmt.Sprintf("https://%s/users/%s/following", hostURL, username),
		Followers:         fmt.Sprintf("https://%s/users/%s/followers", hostURL, username),
		PreferredUsername: username,
		Name:              acc.Name,
		Summary:           acc.Summary,
		URL:               fmt.Sprintf("https://%s/users/%s/hallotest", hostURL, username),
		Published:         acc.CreatedAt.Format(time.RFC3339),
		Inbox:             fmt.Sprintf("https://%s/users/%s/inbox", hostURL, username),
		PublicKey: PublicKey{
			ID:           fmt.Sprintf("https://%s/users/%s#main-key", hostURL, username),
			Owner:        fmt.Sprintf("https://%s/users/%s", hostURL, username),
			PublicKeyPem: string(acc.PublicKey),
		},
	}

}
