package api

import (
	"net/http"

	"github.com/jo-fr/activityhub/modules/api/internal/render"
	"github.com/jo-fr/activityhub/pkg/externalmodel"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
)

func (a *API) AddNewFeedSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req, err := httputil.UnmarshalBody[externalmodel.AddFeedSourceRequest](r.Body)
		if err != nil {
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
