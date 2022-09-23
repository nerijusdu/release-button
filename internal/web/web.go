package web

import (
	"nerijusdu/release-button/internal/argoApi"
	"nerijusdu/release-button/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebApi struct {
	argoApi argoApi.IArgoApi
	config  *config.Config
}

func NewWebApi(aApi argoApi.IArgoApi, c *config.Config) *WebApi {
	return &WebApi{
		argoApi: aApi,
		config:  c,
	}
}

func (a *WebApi) Listen() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	apiRouter := a.getApiRouter()

	r.Mount("/api", apiRouter)

	http.ListenAndServe(":6970", r)
}
