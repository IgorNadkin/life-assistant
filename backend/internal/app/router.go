package app

import (
	"backend/internal/api/handler"

	"github.com/go-chi/chi/v5"
)

func (a *App) initRoutes() {

	healthHandler := handler.NewHealthHandler()
	dbHealth := handler.NewDBHealthHandler(a.db)
	scenarioHandler := handler.NewScenarioHandler(a.scenarioRepo)
	engineHandler := handler.NewEngineHandler(a.engine, a.stateRepo)

	a.router.Route("/api", func(r chi.Router) {

		r.Get("/health", healthHandler.Check)
		r.Get("/health/db", dbHealth.Check)
		r.Post("/scenario", scenarioHandler.Create)
		r.Get("/scenario", scenarioHandler.Get)
		r.Get("/engine/next", engineHandler.Next)
	})
}
