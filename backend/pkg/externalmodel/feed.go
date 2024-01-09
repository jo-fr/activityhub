package externalmodel

import "github.com/jo-fr/activityhub/backend/modules/feed/model"

type AddFeedRequest struct {
	FeedURL string `json:"feedURL" validate:"required,url"`
}

type ListSourcesFeedResponse struct {
	Total int          `json:"total"`
	Items []model.Feed `json:"items"`
}

type ListFeedStatusResponse struct {
	Total int            `json:"total"`
	Items []model.Status `json:"items"`
}
