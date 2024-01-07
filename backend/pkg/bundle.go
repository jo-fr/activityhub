package pkg

import (
	"github.com/jo-fr/activityhub/backend/pkg/config"
	"github.com/jo-fr/activityhub/backend/pkg/database"
	"github.com/jo-fr/activityhub/backend/pkg/log"
	"github.com/jo-fr/activityhub/backend/pkg/pubsub"
	"go.uber.org/fx"
)

var Bundle = fx.Options(
	config.Module,
	log.Module,
	database.Module,
	pubsub.Module,
)
