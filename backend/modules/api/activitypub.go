package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jo-fr/activityhub/backend/modules/api/internal/render"
	"github.com/jo-fr/activityhub/backend/pkg/errutil"
	"github.com/jo-fr/activityhub/backend/pkg/externalmodel"
	"github.com/jo-fr/activityhub/backend/pkg/pubsub"
	"github.com/jo-fr/activityhub/backend/pkg/util/httputil"
	"github.com/jo-fr/activityhub/backend/pkg/validate"
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

		render.Success(r.Context(), externalmodel.ExternalWebfinger(a.hostURL, resource, actor), http.StatusOK, w, a.log)
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

		render.Success(r.Context(), externalmodel.ExternalActor(a.hostURL, actor), http.StatusOK, w, a.log)
	}

}

func (a *API) ReceiveActivity() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		activity, err := httputil.UnmarshalBody[externalmodel.Activity](r.Body)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		if err := validate.Validator().Struct(activity); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		if err := a.pubsub.Publish(r.Context(), pubsub.TopicInbox, activity); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), nil, http.StatusAccepted, w, a.log)
	}

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

		// always returns empty collection because following is not supported
		following := externalmodel.OrderedCollection{
			Context:      "https://www.w3.org/ns/activitystreams",
			ID:           fmt.Sprintf("https://%s/users/%s/following", a.hostURL, actorName),
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

		collection := externalmodel.ExternalFollowerCollection(a.hostURL, actorName, followers)
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
