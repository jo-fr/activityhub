package api

import (
	"context"
	"fmt"

	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/jo-fr/activityhub/backend/modules/activitypub"
	"github.com/jo-fr/activityhub/backend/modules/api/internal/middleware"
	"github.com/jo-fr/activityhub/backend/modules/feed"
	"github.com/jo-fr/activityhub/backend/pkg/config"
	"github.com/jo-fr/activityhub/backend/pkg/log"
	"github.com/jo-fr/activityhub/backend/pkg/pubsub"

	"go.uber.org/fx"
)

const (
	offsetDefault = 0
	limitDefault  = 100
)

var Module = fx.Options(
	fx.Invoke(ProvideAPI),
)

type API struct {
	*chi.Mux
	log     *log.Logger
	host    string
	appHost string

	pubsub      *pubsub.Client
	activitypub *activitypub.Handler
	feed        *feed.Handler
}

func ProvideAPI(lc fx.Lifecycle, config config.Config, logger *log.Logger, pubsub *pubsub.Client, activitypub *activitypub.Handler, feed *feed.Handler) *API {

	api := &API{
		Mux:         chi.NewRouter(),
		log:         logger,
		host:        config.Host,
		appHost:     config.AppHost,
		pubsub:      pubsub,
		activitypub: activitypub,
		feed:        feed,
	}

	api.registerMiddlewares(logger, config)
	api.registerRoutes()

	registerHooks(lc, api, logger, config.Port)

	return api
}

// registerHooks for uber fx
func registerHooks(lc fx.Lifecycle, api *API, logger *log.Logger, port string) {

	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: api}

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

func (a *API) registerMiddlewares(l *log.Logger, config config.Config) {
	a.Use(chiMiddleware.RequestID)
	a.Use(middleware.Logger(l))
	a.Use(chiMiddleware.Recoverer)
	a.Use(middleware.CORSHandler(config.AppHost))

	// add default header
	a.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

}

func (a *API) registerRoutes() {
	a.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK")) // nolint:errcheck
	})

	a.Get("/robots.txt", a.ServeRobotstxt())
	a.Get("/.well-known/webfinger", a.getWebfinger())

	// activitypub relevant endpoints
	a.Route("/ap", func(r chi.Router) {
		r.Get("/{actorName}", a.getActor())
		r.Get("/{actorName}/following", a.FollowingEndpoint())
		r.Get("/{actorName}/followers", a.FollowersEndpoint())

		// protected routes that need a signature header
		r.Group(func(r chi.Router) {
			r.Use(middleware.ValidateSignature(a.log))
			r.Post("/{actorName}/inbox", a.ReceiveActivity())
		})
	})

	a.Route("/api", func(r chi.Router) {
		r.Route("/feeds", func(r chi.Router) {
			r.Post("/", a.AddNewFeed())
			r.Get("/", a.ListFeeds())
			r.Get("/{id}", a.GetFeed())
			r.Get("/{id}/status", a.ListFeedStatus())
		})

		r.Get("/users/{username}/feed", a.GetFeedWithUsername())
		r.Get("/users/{username}/redirect", a.RedirectUser())

	})

}

func (a *API) ServeRobotstxt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := `
# See http://www.robotstxt.org/robotstxt.html for documentation on how to use the robots.txt file

User-agent: GPTBot
Disallow: /

User-agent: *
Disallow: /media_proxy/
Disallow: /interact/
`

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(content)) // nolint:errcheck
	}
}
