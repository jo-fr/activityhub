package externalmodel_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
	"github.com/jo-fr/activityhub/backend/pkg/externalmodel"
	"github.com/jo-fr/activityhub/backend/pkg/sharedmodel"
)

func TestExternalActor(t *testing.T) {
	host := "example.com"
	appHost := "https://app.example.com"
	acc := models.Account{
		BaseModel: sharedmodel.BaseModel{
			ID:        "e8f25aca-b808-45c9-bbe6-08989844ff8e",
			CreatedAt: time.Now(),
		},
		PreferredUsername: "john_doe",
		Name:              "John Doe",
		Summary:           "Lorem ipsum dolor sit amet",
		PublicKey:         []byte("public_key"),
	}

	expectedActor := externalmodel.Actor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		ID:                fmt.Sprintf("https://%s/ap/%s", host, acc.PreferredUsername),
		Type:              "Service",
		Following:         fmt.Sprintf("https://%s/ap/%s/following", host, acc.PreferredUsername),
		Followers:         fmt.Sprintf("https://%s/ap/%s/followers", host, acc.PreferredUsername),
		PreferredUsername: acc.PreferredUsername,
		Name:              acc.Name,
		Summary:           acc.Summary,
		Published:         acc.CreatedAt.Format(time.RFC3339),
		Inbox:             fmt.Sprintf("https://%s/ap/%s/inbox", host, acc.PreferredUsername),
		PublicKey: externalmodel.PublicKey{
			ID:           fmt.Sprintf("https://%s/ap/%s#main-key", host, acc.PreferredUsername),
			Owner:        fmt.Sprintf("https://%s/ap/%s", host, acc.PreferredUsername),
			PublicKeyPem: string(acc.PublicKey),
		},
	}

	result := externalmodel.ExternalActor(host, appHost, acc)

	if !reflect.DeepEqual(result, expectedActor) {
		t.Errorf("Expected %+v, got %+v", expectedActor, result)
	}
}
