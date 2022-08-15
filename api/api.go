package api

import (
	"github.com/go-chi/chi/v5"
)

type Api struct {
	*chi.Mux
}

func NewApi() *Api {
	a := &Api{
		chi.NewRouter(),
	}
	a.Get("/users", a.GetUsers)
	a.Post("/users", a.CreateUser)
	return a
}
