package activitypub

import (
	"fmt"
	"log"

	"github.com/jo-fr/activityhub/modules/activitypub/internal/keys"
	"github.com/jo-fr/activityhub/modules/activitypub/internal/models"
)

var cert = func() *models.KeyPair {
	keys, err := keys.GenerateRSAKeyPair(2048)
	if err != nil {
		log.Fatal(err)
	}
	return keys
}

func (h *Handler) GetActor(actor string) (models.Actor, error) {

	a := models.Actor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		ID:                fmt.Sprintf("%s/%s", h.hostURL, actor),
		Type:              "Person",
		PreferredUsername: actor,
		Inbox:             fmt.Sprintf("%s/inbox", h.hostURL),
		PublicKey: models.PublicKey{
			ID:           fmt.Sprintf("%s/%s#main-key", h.hostURL, actor),
			Owner:        fmt.Sprintf("%s/%s", h.hostURL, actor),
			PublicKeyPem: string(cert().PubKeyPEM),
		},
	}

	return a, nil
}
