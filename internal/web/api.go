package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *WebApi) getApiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/applications", func(w http.ResponseWriter, r *http.Request) {
		apps, err := a.argoApi.GetApps(a.config.Selectors, false)
		if err != nil {
			somethingWentWrong(w, err)
			return
		}

		appsJson, err := json.Marshal(apps)
		if err != nil {
			somethingWentWrong(w, err)
		}
		w.Write(appsJson)
	})

	return r
}

func somethingWentWrong(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("Something went wrong. %v", err)))
}
