package app

import (
	"net/http"

	"backend/internal/api/handler"
	"backend/internal/config"
	"backend/internal/repository/postgres"
	"backend/internal/service"
	"backend/internal/service/engine"
	"backend/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type App struct {
	router *chi.Mux
	cfg    config.Config
	db     *sqlx.DB
}

func New() (*App, error) {
	cfg := config.Load()
	db := database.NewPostgres(cfg.PostgresDSN)

	// Initialize repositories
	scenarioRepo := postgres.NewScenarioRepo(db)
	nodeRepo := postgres.NewNodeRepo(db)
	edgeRepo := postgres.NewEdgeRepo(db)
	userRepo := postgres.NewUserRepo(db)
	userStateRepo := postgres.NewUserStateRepo(db)
	userActionRepo := postgres.NewUserActionRepo(db)

	// Initialize services
	scenarioService := service.NewScenarioService(scenarioRepo, userStateRepo, nodeRepo)
	nodeService := service.NewNodeService(nodeRepo)
	edgeService := service.NewEdgeService(edgeRepo)
	userService := service.NewUserService(userRepo, userStateRepo, userActionRepo)

	// Initialize engine
	nodes, _ := nodeRepo.GetByScenario(0) // TODO: load nodes by scenario
	edges, _ := edgeRepo.GetByScenario(0) // TODO: load edges by scenario
	graphEngine := engine.New(nodes, edges)
	scenarioFlow := engine.NewScenarioFlow(graphEngine, userStateRepo, userActionRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	dbHealthHandler := handler.NewDBHealthHandler(db)
	scenarioHandler := handler.NewScenarioHandler(scenarioService)
	nodeHandler := handler.NewNodeHandler(nodeService)
	edgeHandler := handler.NewEdgeHandler(edgeService)
	engineHandler := handler.NewEngineHandler(scenarioService, userService, scenarioFlow, userStateRepo)

	r := chi.NewRouter()

	app := &App{
		router: r,
		cfg:    cfg,
		db:     db,
	}

	app.initRoutes(healthHandler, dbHealthHandler, scenarioHandler, nodeHandler, edgeHandler, engineHandler)

	return app, nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":" + a.cfg.Port, a.router)
}
