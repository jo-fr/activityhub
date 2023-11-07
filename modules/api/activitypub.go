package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	model "github.com/jo-fr/activityhub/modules/api/internal/externalmodel"
	"github.com/jo-fr/activityhub/modules/api/internal/render"
	"github.com/jo-fr/activityhub/pkg/errutil"
)

// api errors
var (
	ErrWrongURI    = errutil.NewError(errutil.TypeInvalidRequestBody, "uri of actor and host do not match")
	ErrWrongFormat = errutil.NewError(errutil.TypeBadRequest, "no acct: in resource")
)

func (a *API) getWebfinger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resource := r.URL.Query().Get("resource")

		username, err := validateAndExtractUsername(a.hostURL, resource)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		actor, err := a.activitypub.GetActor(username)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), model.ExternalWebfinger(a.hostURL, resource, actor), http.StatusOK, w, a.log)
	}
}

func (a *API) getActor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		actorName := chi.URLParam(r, "actorName")

		actor, err := a.activitypub.GetActor(actorName)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), model.ExternalActor(a.hostURL, actor), http.StatusOK, w, a.log)
	}

}

func validateAndExtractUsername(hostURL string, resource string) (string, error) {
	if !strings.Contains(resource, "acct:") {
		return "", ErrWrongFormat
	}

	actor := strings.Split(resource, "acct:")[1]
	username := strings.Split(actor, "@")[0]
	uri := strings.Split(actor, "@")[1]

	if uri != hostURL {
		return "", ErrWrongURI
	}

	return username, nil
}
