package api

import (
	"bytes"
	"context"
	"io"

	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/modules/api/internal/middleware"
	"github.com/jo-fr/activityhub/modules/api/internal/render"
	"github.com/jo-fr/activityhub/pkg/config"
	"github.com/jo-fr/activityhub/pkg/log"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(ProvideAPI),
)

type API struct {
	*chi.Mux
	log     *log.Logger
	hostURL string

	activitypub *activitypub.Handler
}

func ProvideAPI(lc fx.Lifecycle, config config.Config, logger *log.Logger, activitypub *activitypub.Handler) *API {

	api := &API{
		Mux:         chi.NewRouter(),
		log:         logger,
		hostURL:     config.HostURL,
		activitypub: activitypub,
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

	a.Get("/.well-known/webfinger", a.getWebfinger())
	a.Get("/{actorName}", a.getActor())

	a.Post("/inbox", a.inbox())
	a.Get("/inbox", a.inbox())

}

func (a *API) inbox() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		actorName := chi.URLParam(r, "actorName")
		a.log.Info("Actor name: " + actorName)

		body, err := readBodyToString(r.Body)
		if err != nil {
			a.log.Fatal(err)
		}

		a.log.Info("Request body: " + body)

		actor, err := a.activitypub.GetActor(actorName)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), actor, http.StatusOK, w, a.log)

	}
}
func readBodyToString(body io.ReadCloser) (string, error) {
	defer body.Close()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
