package activitypub

import (
	"github.com/jo-fr/activityhub/modules/activitypub/models"
)

func (h *Handler) GetActor(actor string) (models.Account, error) {

	var acc models.Account
	if err := h.db.First(&acc, "preferred_username = ?", actor).Error; err != nil {
		return models.Account{}, err
	}

	return acc, nil
}
