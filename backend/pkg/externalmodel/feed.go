package externalmodel

import (
	"github.com/jo-fr/activityhub/backend/modules/feed/model"
)

type AddFeedRequest struct {
	FeedURL string `json:"feedURL" validate:"required,url"`
}

type ListSourcesFeedResponse struct {
	Total int    `json:"total"`
	Items []Feed `json:"items"`
}

type ListFeedStatusResponse struct {
	Total int            `json:"total"`
	Items []model.Status `json:"items"`
}

type Feed struct {
	CreatedAt   string         `json:"createdAt"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        model.FeedType `json:"type"`
	FeedURL     string         `json:"feedURL"`
	HostURL     string         `json:"host"`
	Author      string         `json:"author"`
	Description string         `json:"description"`
	ImageURL    string         `json:"imageURL"`
	Account     Account        `json:"account"`
}

func ExternalFeed(feed model.Feed, host string) Feed {
	return Feed{
		CreatedAt:   feed.CreatedAt.String(),
		ID:          feed.ID,
		Name:        feed.Name,
		Type:        feed.Type,
		FeedURL:     feed.FeedURL,
		HostURL:     feed.HostURL,
		Author:      feed.Author,
		Description: feed.Description,
		ImageURL:    feed.ImageURL,
		Account:     ExternalAccount(feed.Account, host),
	}

}
