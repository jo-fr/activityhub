package externalmodel_test

import (
	"fmt"
	"testing"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
	"github.com/jo-fr/activityhub/backend/pkg/externalmodel"
)

func TestExternalWebfinger(t *testing.T) {
	host := "example.com"
	appHost := "app.example.com"
	resource := "user123"
	acc := models.Account{
		PreferredUsername: "user123",
	}

	expectedWebfinger := externalmodel.Webfinger{
		Subject: resource,
		Links: []externalmodel.Links{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s/ap/%s", host, acc.PreferredUsername),
			},
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: fmt.Sprintf("https://%s/feed/%s", appHost, acc.PreferredUsername),
			},
		},
	}

	result := externalmodel.ExternalWebfinger(host, appHost, resource, acc)

	if result.Subject != expectedWebfinger.Subject {
		t.Errorf("Expected Subject %s, got %s", expectedWebfinger.Subject, result.Subject)
	}

	if len(result.Links) != len(expectedWebfinger.Links) {
		t.Errorf("Expected %d links, got %d", len(expectedWebfinger.Links), len(result.Links))
	}

	for i := range result.Links {
		if result.Links[i].Rel != expectedWebfinger.Links[i].Rel {
			t.Errorf("Expected Rel %s, got %s", expectedWebfinger.Links[i].Rel, result.Links[i].Rel)
		}

		if result.Links[i].Type != expectedWebfinger.Links[i].Type {
			t.Errorf("Expected Type %s, got %s", expectedWebfinger.Links[i].Type, result.Links[i].Type)
		}

		if result.Links[i].Href != expectedWebfinger.Links[i].Href {
			t.Errorf("Expected Href %s, got %s", expectedWebfinger.Links[i].Href, result.Links[i].Href)
		}
	}
}
