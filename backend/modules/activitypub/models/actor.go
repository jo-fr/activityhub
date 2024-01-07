package models

import "github.com/jo-fr/activityhub/backend/pkg/sharedmodel"

type Account struct {
	sharedmodel.BaseModel
	PreferredUsername string
	Name              string
	Summary           string
	PrivateKey        []byte
	PublicKey         []byte
}

func (*Account) TableName() string {
	return "activityhub.account"
}
