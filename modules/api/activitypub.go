package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	model "github.com/jo-fr/activityhub/modules/api/externalmodel"
	"github.com/jo-fr/activityhub/modules/api/internal/render"
	"github.com/jo-fr/activityhub/pkg/errutil"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
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

		actor, err := a.activitypub.GetActor(r.Context(), username)
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

		actor, err := a.activitypub.GetActor(r.Context(), actorName)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), model.ExternalActor(a.hostURL, actor), http.StatusOK, w, a.log)
	}

}

func (a *API) ReceiveActivity() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		activity, err := httputil.UnmarshaBody[model.Activity](r.Body)
		if err != nil {
			render.Success(r.Context(), nil, http.StatusOK, w, a.log)
			return
		}

		a.log.Info(activity)

		err = a.activitypub.ReceiveInboxActivity(r.Context(), activity)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), nil, http.StatusOK, w, a.log)
	}

}

type OrderedCollection struct {
	Context      string   `json:"@context"`
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	TotalItems   int      `json:"totalItems"`
	OrderedItems []string `json:"orderedItems"`
}

func (a *API) FollowingEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		actorName := chi.URLParam(r, "actorName")
		// check if actor exists on instance
		_, err := a.activitypub.GetActor(r.Context(), actorName)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		following := OrderedCollection{
			Context:      "https://www.w3.org/ns/activitystreams",
			ID:           fmt.Sprintf("https://%s/users/joni/following", a.hostURL),
			Type:         "OrderedCollection",
			TotalItems:   0,
			OrderedItems: []string{},
		}

		render.Success(r.Context(), following, http.StatusOK, w, a.log)
	}
}

func (a *API) FollowersEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		actorName := chi.URLParam(r, "actorName")
		followers, err := a.activitypub.GetFollowers(r.Context(), actorName)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		collection := model.ExternalFollowerCollection(a.hostURL, actorName, followers)
		render.Success(r.Context(), collection, http.StatusOK, w, a.log)
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
