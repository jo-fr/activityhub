package models

import "github.com/jo-fr/activityhub/backend/pkg/sharedmodel"

type Account struct {
	sharedmodel.BaseModel `json:"-"`
	PreferredUsername     string `json:"preferredUsername"`
	Name                  string `json:"name"`
	Summary               string `json:"summary"`
	PrivateKey            []byte `json:"-"`
	PublicKey             []byte `json:"-"`
}

func (*Account) TableName() string {
	return "activityhub.account"
}
