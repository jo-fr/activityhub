package pkg

import (
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/jo-fr/activityhub/pkg/database"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/jo-fr/activityhub/pkg/pubsub"
	"go.uber.org/fx"
)

var Bundle = fx.Options(
	config.Module,
	log.Module,
	database.Module,
	pubsub.Module,
)
