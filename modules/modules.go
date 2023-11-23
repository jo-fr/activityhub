package modules

import (
	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/modules/api"
	"github.com/jo-fr/activityhub/modules/feed"
	"go.uber.org/fx"
)

var Bundle = fx.Options(
	feed.Module,
	activitypub.Module,
	api.Module,
)
