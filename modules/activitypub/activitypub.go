package activitypub

import (
	"log"

	"github.com/jo-fr/activityhub/modules/activitypub/internal/keys"
	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/jo-fr/activityhub/pkg/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(ProvideHandler),
)

type Handler struct {
	hostURL string
	db      *database.Database
}

func ProvideHandler(config config.Config, db *database.Database) *Handler {

	pair, err := keys.GenerateRSAKeyPair(2048)
	if err != nil {
		log.Fatalln(err, "failed to generate RSA key pair")
	}

	if err := db.Create(&models.Account{
		PreferredUsername: "joni",
		Name:              "Jonathan",
		Summary:           "This is the profile of Jonathan",
		PrivateKey:        pair.PrivKeyPEM,
		PublicKey:         pair.PubKeyPEM,
	}).Error; err != nil {
		log.Fatalln(err, "failed to generate account key pair")
	}

	return &Handler{
		hostURL: config.HostURL,
		db:      db,
	}
}
