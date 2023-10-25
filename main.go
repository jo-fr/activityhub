package main

import (
	"github.com/jo-fr/activityhub/modules"
	"github.com/jo-fr/activityhub/pkg"
	"go.uber.org/fx"
)

var App = fx.New(
	pkg.Bundle,
	modules.Bundle,
)

func main() {
	App.Run()
}
