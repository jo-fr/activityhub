package modules

import (
	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/modules/api"
	"go.uber.org/fx"
)

var Bundle = fx.Options(
	activitypub.Module,
	api.Module,
)
