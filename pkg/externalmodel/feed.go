package externalmodel

import "github.com/jo-fr/activityhub/modules/feed/model"

type AddFeedSourceRequest struct {
	FeedURL string `json:"feedURL"`
}

type ListSourcesFeedResponse struct {
	Total int                `json:"total"`
	Items []model.SourceFeed `json:"items"`
}
