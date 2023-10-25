package api

import (
	"context"

	"net/http"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/pkg/log"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(ProvideAPI),
)

type API struct {
	*chi.Mux
	log *log.Logger

	activitypub *activitypub.Handler
}

func ProvideAPI(lc fx.Lifecycle, logger *log.Logger, activitypub *activitypub.Handler) *API {

	api := &API{
		Mux:         chi.NewRouter(),
		log:         logger,
		activitypub: activitypub,
	}

	api.registerMiddlewares()
	api.registerRoutes()

	registerHooks(lc, api, logger)

	return api
}

// registerHooks for uber fx
func registerHooks(lc fx.Lifecycle, api *API, logger *log.Logger) {

	server := &http.Server{Addr: ":8080", Handler: api}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("starting api server")
				go func() {
					err := server.ListenAndServe()
					if err != nil && err != http.ErrServerClosed {
						logger.Fatal(err)
					}

				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("shutting down api server")
				return server.Shutdown(ctx)
			},
		},
	)
}

func (a *API) registerMiddlewares() {
	a.Use(chiMiddleware.RequestID)
	a.Use(chiMiddleware.Logger)
	a.Use(chiMiddleware.Recoverer)
	// add default header
	a.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

}

func (a *API) registerRoutes() {
	a.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	a.Get("/.well-known/webfinger", a.getWebfinger())

}
