package model

import "github.com/jo-fr/activityhub/pkg/sharedmodel"

type FeedType string

const (
	FeedTypeRSS FeedType = "RSS"
)

type Feed struct {
	sharedmodel.BaseModel
	Name        string   `json:"name"`
	Type        FeedType `json:"type"`
	FeedURL     string   `json:"feedURL"`
	HostURL     string   `json:"hostURL"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	ImageURL    string   `json:"imageURL"`
	AccountID   string   `json:"accountID"`
}

func (*Feed) TableName() string {
	return "activityhub.feed"
}
