package modules

import (
	"github.com/jo-fr/activityhub/backend/modules/activitypub"
	"github.com/jo-fr/activityhub/backend/modules/api"
	"github.com/jo-fr/activityhub/backend/modules/feed"
	"go.uber.org/fx"
)

var Bundle = fx.Options(
	feed.Module,
	activitypub.Module,
	api.Module,
)
