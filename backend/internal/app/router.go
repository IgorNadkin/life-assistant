package app

import (
	"life-assistant/internal/api/handler"

	"github.com/go-chi/chi/v5"
)

func (a *App) initRoutes() {

	healthHandler := handler.NewHealthHandler()

	a.router.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.Check)
	})
}
