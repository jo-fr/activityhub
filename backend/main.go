package main

import (
	"github.com/jo-fr/activityhub/backend/modules"
	"github.com/jo-fr/activityhub/backend/pkg"
	"go.uber.org/fx"
)

var App = fx.New(
	pkg.Bundle,
	modules.Bundle,
)

func main() {
	App.Run()
}
