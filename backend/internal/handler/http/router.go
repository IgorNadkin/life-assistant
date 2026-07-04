package http

import (
	"life-assistant/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(s *service.Services) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {

		r.Post("/session/start", s.Session.Start)
		r.Post("/session/next", s.Session.Next)

		r.Post("/scenario", s.Scenario.Create)
		r.Get("/scenario/{id}", s.Scenario.Get)
	})

	return r
}
