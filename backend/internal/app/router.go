package app

import (
	"backend/internal/api/handler"

	"github.com/go-chi/chi/v5"
)

func (a *App) initRoutes(
	healthHandler *handler.HealthHandler,
	dbHealthHandler *handler.DBHealthHandler,
	scenarioHandler *handler.ScenarioHandler,
	nodeHandler *handler.NodeHandler,
	edgeHandler *handler.EdgeHandler,
	engineHandler *handler.EngineHandler,
) {
	a.router.Route("/api", func(r chi.Router) {
		// Health checks
		r.Get("/health", healthHandler.Check)
		r.Get("/health/db", dbHealthHandler.Check)

		// Scenarios
		r.Post("/scenario", scenarioHandler.Create)
		r.Get("/scenario", scenarioHandler.Get)
		r.Put("/scenario", scenarioHandler.Update)
		r.Delete("/scenario", scenarioHandler.Delete)
		r.Get("/scenarios", scenarioHandler.List)

		// Nodes
		r.Post("/node", nodeHandler.Create)
		r.Get("/node", nodeHandler.Get)
		r.Get("/nodes", nodeHandler.GetByScenario)
		r.Put("/node", nodeHandler.Update)
		r.Delete("/node", nodeHandler.Delete)

		// Edges
		r.Post("/edge", edgeHandler.Create)
		r.Get("/edge", edgeHandler.Get)
		r.Get("/edges", edgeHandler.GetByScenario)
		r.Put("/edge", edgeHandler.Update)
		r.Delete("/edge", edgeHandler.Delete)

		// Scenario Engine
		r.Post("/user/scenario/start", engineHandler.Start)
		r.Post("/user/scenario/step", engineHandler.Step)
		r.Get("/user/scenario/status", engineHandler.GetStatus)
	})
}
