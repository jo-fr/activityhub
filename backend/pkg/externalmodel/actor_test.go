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
	hostURL := "example.com"
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
		ID:                fmt.Sprintf("https://%s/users/%s", hostURL, acc.PreferredUsername),
		Type:              "Service",
		Following:         fmt.Sprintf("https://%s/users/%s/following", hostURL, acc.PreferredUsername),
		Followers:         fmt.Sprintf("https://%s/users/%s/followers", hostURL, acc.PreferredUsername),
		PreferredUsername: acc.PreferredUsername,
		Name:              acc.Name,
		Summary:           acc.Summary,
		Published:         acc.CreatedAt.Format(time.RFC3339),
		Inbox:             fmt.Sprintf("https://%s/users/%s/inbox", hostURL, acc.PreferredUsername),
		PublicKey: externalmodel.PublicKey{
			ID:           fmt.Sprintf("https://%s/users/%s#main-key", hostURL, acc.PreferredUsername),
			Owner:        fmt.Sprintf("https://%s/users/%s", hostURL, acc.PreferredUsername),
			PublicKeyPem: string(acc.PublicKey),
		},
	}

	result := externalmodel.ExternalActor(hostURL, acc)

	if !reflect.DeepEqual(result, expectedActor) {
		t.Errorf("Expected %+v, got %+v", expectedActor, result)
	}
}
