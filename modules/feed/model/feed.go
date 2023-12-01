package model

import "github.com/jo-fr/activityhub/pkg/sharedmodel"

type SourceFeedType string

const (
	SourceFeedTypeRSS SourceFeedType = "RSS"
)

type SourceFeed struct {
	sharedmodel.BaseModel
	Name        string
	Type        SourceFeedType
	FeedURL     string
	HostURL     string
	Author      string
	Description string
	ImageURL    string
	AccountID   string
}

func (*SourceFeed) TableName() string {
	return "activityhub.source_feed"
}
