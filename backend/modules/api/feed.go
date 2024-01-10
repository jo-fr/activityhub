package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jo-fr/activityhub/backend/modules/api/internal/render"
	"github.com/jo-fr/activityhub/backend/pkg/externalmodel"
	"github.com/jo-fr/activityhub/backend/pkg/util/httputil"
	"github.com/jo-fr/activityhub/backend/pkg/validate"
)

func (a *API) AddNewFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req, err := httputil.UnmarshalRequestBody[externalmodel.AddFeedRequest](r)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		if err := validate.Validator().Struct(req); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		feed, err := a.feed.AddNewFeed(r.Context(), req.FeedURL)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), externalmodel.ExternalFeed(feed, a.host), http.StatusCreated, w, a.log)
	}
}

func (a *API) ListFeeds() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		offset, err := strconv.ParseInt(r.FormValue("offset"), 10, 32)
		if err != nil || offset == 0 {
			offset = offsetDefault
		}

		limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
		if err != nil || limit == 0 {
			limit = limitDefault
		}

		totalCount, sources, err := a.feed.ListFeeds(r.Context(), int(offset), int(limit))
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		var extFeeds []externalmodel.Feed
		for _, feed := range sources {
			extFeeds = append(extFeeds, externalmodel.ExternalFeed(feed, a.host))
		}

		resp := externalmodel.ListSourcesFeedResponse{
			Total: totalCount,
			Items: extFeeds,
		}

		render.Success(r.Context(), resp, http.StatusOK, w, a.log)
	}
}

func (a *API) GetFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		if err := validate.Validator().Var(id, "required,uuid4"); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		feed, err := a.feed.GetFeed(r.Context(), id)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), externalmodel.ExternalFeed(feed, a.host), http.StatusOK, w, a.log)
	}
}

func (a *API) GetFeedWithUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username := chi.URLParam(r, "username")

		feed, err := a.feed.GetFeedWithUsername(r.Context(), username)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), externalmodel.ExternalFeed(feed, a.host), http.StatusOK, w, a.log)
	}
}

// Redirect user is a workaround to link from a mastodoan instance to the web app
func (a *API) RedirectUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username := chi.URLParam(r, "username")
		feedURL := fmt.Sprintf("https://%s/feed/%s", a.appHost, username)

		http.Redirect(w, r, feedURL, http.StatusPermanentRedirect)
	}
}

func (a *API) ListFeedStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		if err := validate.Validator().Var(id, "required,uuid4"); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		offset, err := strconv.ParseInt(r.FormValue("offset"), 10, 32)
		if err != nil || offset == 0 {
			offset = offsetDefault
		}

		limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
		if err != nil || limit == 0 {
			limit = limitDefault
		}

		totalCount, status, err := a.feed.ListFeedStatus(r.Context(), id, int(offset), int(limit))
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		resp := externalmodel.ListFeedStatusResponse{
			Total: totalCount,
			Items: status,
		}

		render.Success(r.Context(), resp, http.StatusOK, w, a.log)
	}
}
