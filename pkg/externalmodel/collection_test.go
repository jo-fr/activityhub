package externalmodel_test

import (
	"fmt"
	"testing"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/jo-fr/activityhub/pkg/externalmodel"
)

func TestExternalFollowerCollection(t *testing.T) {
	hostURL := "example.com"
	username := "john"
	followers := []models.Follower{
		{AccountURIFollowing: "https://example.com/user1"},
		{AccountURIFollowing: "https://example.com/user2"},
		{AccountURIFollowing: "https://example.com/user3"},
	}

	expectedCollection := externalmodel.OrderedCollection{
		Context:      "https://www.w3.org/ns/activitystreams",
		ID:           fmt.Sprintf("https://%s/users/%s/followers", hostURL, username),
		Type:         "OrderedCollection",
		TotalItems:   len(followers),
		OrderedItems: []string{"https://example.com/user1", "https://example.com/user2", "https://example.com/user3"},
	}

	result := externalmodel.ExternalFollowerCollection(hostURL, username, followers)

	if result.Context != expectedCollection.Context {
		t.Errorf("Expected Context %s, got %s", expectedCollection.Context, result.Context)
	}

	if result.ID != expectedCollection.ID {
		t.Errorf("Expected ID %s, got %s", expectedCollection.ID, result.ID)
	}

	if result.Type != expectedCollection.Type {
		t.Errorf("Expected Type %s, got %s", expectedCollection.Type, result.Type)
	}

	if result.TotalItems != expectedCollection.TotalItems {
		t.Errorf("Expected TotalItems %d, got %d", expectedCollection.TotalItems, result.TotalItems)
	}

	if len(result.OrderedItems) != len(expectedCollection.OrderedItems) {
		t.Errorf("Expected OrderedItems length %d, got %d", len(expectedCollection.OrderedItems), len(result.OrderedItems))
	} else {
		for i := range result.OrderedItems {
			if result.OrderedItems[i] != expectedCollection.OrderedItems[i] {
				t.Errorf("Expected OrderedItems[%d] %s, got %s", i, expectedCollection.OrderedItems[i], result.OrderedItems[i])
			}
		}
	}
}
