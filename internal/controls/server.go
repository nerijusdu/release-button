package controls

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type IOServer struct {
}

func NewIOServer() *IOServer {
	return &IOServer{}
}

func (c *IOServer) Listen(clickChan chan<- string) {
	fmt.Println("Listening to IO requests")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/io/buttons/{button}", func(w http.ResponseWriter, r *http.Request) {
		clickChan <- chi.URLParam(r, "button")
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":6969", r)
}