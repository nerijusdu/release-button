package web

import (
	"fmt"
	"nerijusdu/release-button/internal/argoApi"
	"nerijusdu/release-button/internal/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "/web/dist"))
	FileServer(r, "/", filesDir)

	fmt.Println("Running Web UI on port :6969")
	http.ListenAndServe(":6969", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
