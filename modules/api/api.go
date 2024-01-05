package api

import (
	"context"

	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/modules/api/internal/middleware"
	"github.com/jo-fr/activityhub/modules/feed"
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/jo-fr/activityhub/pkg/pubsub"

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
	hostURL string

	pubsub      *pubsub.Client
	activitypub *activitypub.Handler
	feed        *feed.Handler
}

func ProvideAPI(lc fx.Lifecycle, config config.Config, logger *log.Logger, pubsub *pubsub.Client, activitypub *activitypub.Handler, feed *feed.Handler) *API {

	api := &API{
		Mux:         chi.NewRouter(),
		log:         logger,
		hostURL:     config.HostURL,
		pubsub:      pubsub,
		activitypub: activitypub,
		feed:        feed,
	}

	api.registerMiddlewares(logger)
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

func (a *API) registerMiddlewares(l *log.Logger) {
	a.Use(chiMiddleware.RequestID)
	a.Use(middleware.Logger(l))
	a.Use(chiMiddleware.Recoverer)

	// add default header
	a.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

}

func (a *API) registerRoutes() {
	a.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	a.Get("/robots.txt", a.ServeRobotstxt())
	a.Get("/.well-known/webfinger", a.getWebfinger())

	// /users endpoints
	a.Route("/users", func(r chi.Router) {

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
		r.Post("/feed", a.AddNewFeed())
		r.Get("/feed", a.ListFeeds())
		r.Get("/feed/{id}", a.GetFeed())
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
		w.Write([]byte(content))
	}
}
