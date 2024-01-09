package externalmodel

import (
	"fmt"

	"github.com/jo-fr/activityhub/backend/modules/activitypub/models"
)

type Account struct {
	CreatedAt string `json:"createdAt"`
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
}

func ExternalAccount(acc models.Account, host string) Account {
	return Account{
		CreatedAt: acc.CreatedAt.String(),
		ID:        acc.ID,
		Username:  acc.PreferredUsername,
		Name:      acc.Name,
		URI:       fmt.Sprintf("%s@%s", acc.PreferredUsername, host),
	}
}
