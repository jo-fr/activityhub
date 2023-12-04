package model

import "github.com/jo-fr/activityhub/pkg/sharedmodel"

type Status struct {
	sharedmodel.BaseModel
	Content   string
	AccountID string
}

func (*Status) TableName() string {
	return "activityhub.status"
}
