package api

import (
	"net/http"

	"github.com/jo-fr/activityhub/modules/api/internal/render"
)

func (a *API) getWebfinger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resource := r.URL.Query().Get("resource")

		webfinger, err := a.activitypub.GetWebfinger(resource)
		if err != nil {
			render.Error(r.Context(), err, w, a.log)
			return
		}

		render.Success(r.Context(), webfinger, http.StatusOK, w, a.log)
	}
}
