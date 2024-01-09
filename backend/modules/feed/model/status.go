package model

import "github.com/jo-fr/activityhub/backend/pkg/sharedmodel"

type Status struct {
	sharedmodel.BaseModel
	Content   string `json:"content"`
	AccountID string `json:"accountID"`
}

func (*Status) TableName() string {
	return "activityhub.status"
}
