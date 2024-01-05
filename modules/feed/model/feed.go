package model

import "github.com/jo-fr/activityhub/pkg/sharedmodel"

type SourceFeedType string

const (
	SourceFeedTypeRSS SourceFeedType = "RSS"
)

type SourceFeed struct {
	sharedmodel.BaseModel
	Name        string         `json:"name"`
	Type        SourceFeedType `json:"type"`
	FeedURL     string         `json:"feedURL"`
	HostURL     string         `json:"hostURL"`
	Author      string         `json:"author"`
	Description string         `json:"description"`
	ImageURL    string         `json:"imageURL"`
	AccountID   string         `json:"accountID"`
}

func (*SourceFeed) TableName() string {
	return "activityhub.source_feed"
}
