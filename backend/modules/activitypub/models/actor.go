package models

import "github.com/jo-fr/activityhub/backend/pkg/sharedmodel"

type Account struct {
	sharedmodel.BaseModel `json:"-"`
	PreferredUsername     string `json:"preferredUsername,omitempty"`
	Name                  string `json:"name,omitempty"`
	Summary               string `json:"summary,omitempty"`
	PrivateKey            []byte `json:"-"`
	PublicKey             []byte `json:"-"`
}

func (*Account) TableName() string {
	return "activityhub.account"
}
