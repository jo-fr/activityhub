package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jo-fr/activityhub/modules/api/internal/render"
	"github.com/jo-fr/activityhub/pkg/externalmodel"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
	"github.com/jo-fr/activityhub/pkg/validate"
)

func (a *API) AddNewFeedSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req, err := httputil.UnmarshalBody[externalmodel.AddFeedSourceRequest](r.Body)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		if err := validate.Validator().Struct(req); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		sourceFeed, err := a.feed.AddNewSourceFeed(r.Context(), req.FeedURL)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), sourceFeed, http.StatusCreated, w, a.log)
	}
}

func (a *API) ListFeedSources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		offset, err := strconv.ParseInt(r.FormValue("offset"), 10, 32)
		if err != nil || offset == 0 {
			offset = offsetDefault
		}

		limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
		if err != nil || limit == 0 {
			limit = limitDefault
		}

		totalCount, sources, err := a.feed.ListSourceFeeds(r.Context(), int(offset), int(limit))
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		resp := externalmodel.ListSourcesFeedResponse{
			Total: totalCount,
			Items: sources,
		}

		render.Success(r.Context(), resp, http.StatusOK, w, a.log)
	}
}

func (a *API) GetFeedSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		if err := validate.Validator().Var(id, "required,uuid4"); err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		source, err := a.feed.GetSourceFeed(r.Context(), id)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), source, http.StatusOK, w, a.log)
	}
}
