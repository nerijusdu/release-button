package controls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type IOListener struct {
}

type ActionData struct {
	Number int `json:"number"`
}

type Action struct {
	Action string `json:"action"`
	Data   ActionData
}

func NewIOListener() *IOListener {
	return &IOListener{}
}

func (c *IOListener) Listen(clickChan chan<- Action) {
	fmt.Println("Listening to IO requests")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/io/actions/{action}", func(w http.ResponseWriter, r *http.Request) {
		d := ActionData{}
		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("ERR: failed to read body. %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if len(bodyBytes) != 0 {
			err = json.Unmarshal(bodyBytes, &d)
			if err != nil {
				fmt.Printf("ERR: failed to unmarshal body. %v", err)
				w.WriteHeader(http.StatusBadRequest)
			}
		}

		a := Action{
			Action: chi.URLParam(r, "action"),
			Data:   d,
		}

		clickChan <- a

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Running controls listener on port :6970")
	http.ListenAndServe(":6970", r)
}
