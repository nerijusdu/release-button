package web

import (
	"encoding/json"
	"fmt"
	"nerijusdu/release-button/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *WebApi) getApiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))

	r.Get("/applications", func(w http.ResponseWriter, r *http.Request) {
		apps, err := a.argoApi.GetApps(a.config.Selectors, false)
		if err != nil {
			somethingWentWrong(w, err)
			return
		}

		appsJson, err := json.Marshal(apps)
		if err != nil {
			somethingWentWrong(w, err)
			return
		}

		w.Write(appsJson)
	})

	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		configJson, err := json.Marshal(a.config)
		if err != nil {
			somethingWentWrong(w, err)
			return
		}

		w.Write(configJson)
	})

	r.Post("/config", func(w http.ResponseWriter, r *http.Request) {
		data := config.Config{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			somethingWentWrong(w, err)
			return
		}

		a.config.Allowed = data.Allowed
		a.config.Selectors = data.Selectors
		a.config.RefreshInterval = data.RefreshInterval

		config.WriteConfig(*a.config)

		w.WriteHeader(http.StatusOK)
	})

	return r
}

func somethingWentWrong(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("Something went wrong. %v", err)))
}
