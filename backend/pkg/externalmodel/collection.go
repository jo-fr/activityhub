package externalmodel

import (
	"fmt"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
)

type OrderedCollection struct {
	Context      string   `json:"@context"`
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	TotalItems   int      `json:"totalItems"`
	OrderedItems []string `json:"orderedItems"`
}

func ExternalFollowerCollection(hostURL string, username string, followers []models.Follower) OrderedCollection {
	var orderedItems []string
	for _, follower := range followers {
		orderedItems = append(orderedItems, follower.AccountURIFollowing)
	}
	return OrderedCollection{
		Context:      "https://www.w3.org/ns/activitystreams",
		ID:           fmt.Sprintf("https://%s/ap/%s/followers", hostURL, username),
		Type:         "OrderedCollection",
		TotalItems:   len(followers),
		OrderedItems: orderedItems,
	}
}
