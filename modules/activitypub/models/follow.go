package models

import "github.com/jo-fr/activityhub/pkg/sharedmodel"

type Follower struct {
	sharedmodel.BaseModel
	AccountIDFollowed   string `gorm:"not null"`
	AccountURIFollowing string `gorm:"not null"`
}

func (*Follower) TableName() string {
	return "activityhub.follower"
}
