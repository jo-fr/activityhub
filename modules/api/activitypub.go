package api

import (
	"encoding/json"
	"fmt"
	"io"
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

func (a *API) ReceivceActivity() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		a.log.Info(1)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}
		a.log.Info(2)
		var activity model.Activity
		if err := json.Unmarshal(body, &activity); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		a.log.Info(activity)

		_, err = a.activitypub.ReceiveInboxActivity(r.Context(), activity.Actor, activity.Object, activity.Type)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}
		a.log.Info(4)

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

		following := OrderedCollection{
			Context:    "https://www.w3.org/ns/activitystreams",
			ID:         fmt.Sprintf("https://%s/users/joni/following", a.hostURL),
			Type:       "OrderedCollection",
			TotalItems: 4,
			OrderedItems: []string{
				"https://tldr.nettime.org/users/tante",
				"https://social.hetzel.net/users/timo",
				"https://social.rebellion.global/users/ScientistRebellion",
				"https://social.network.europa.eu/users/EU_Commission",
			},
		}

		render.Success(r.Context(), following, http.StatusOK, w, a.log)
	}
}

func (a *API) FollowersEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		followers := OrderedCollection{
			Context:    "https://www.w3.org/ns/activitystreams",
			ID:         fmt.Sprintf("https://%s/users/joni/followers", a.hostURL),
			Type:       "OrderedCollection",
			TotalItems: 4,
			OrderedItems: []string{
				"https://tldr.nettime.org/users/tante",
				"https://social.hetzel.net/users/timo",
				"https://social.rebellion.global/users/ScientistRebellion",
				"https://social.network.europa.eu/users/EU_Commission",
			},
		}

		render.Success(r.Context(), followers, http.StatusOK, w, a.log)
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
