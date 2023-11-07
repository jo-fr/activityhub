package model

import (
	"fmt"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

type Actor struct {
	Context           []string  `json:"@context"`
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	PreferredUsername string    `json:"preferredUsername"`
	Name              string    `json:"name"`
	Summary           string    `json:"summary"`
	Inbox             string    `json:"inbox"`
	PublicKey         PublicKey `json:"publicKey"`
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
		ID:                fmt.Sprintf("https://%s/%s", hostURL, username),
		Type:              "Person",
		PreferredUsername: username,
		Name:              acc.Name,
		Summary:           acc.Summary,

		Inbox: fmt.Sprintf("https://%s/%s/inbox", hostURL, username),
		PublicKey: PublicKey{
			ID:           fmt.Sprintf("https://%s/%s#main-key", hostURL, username),
			Owner:        fmt.Sprintf("https://%s/%s", hostURL, username),
			PublicKeyPem: string(acc.PublicKey),
		},
	}

}
