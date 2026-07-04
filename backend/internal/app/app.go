package app

import (
	"net/http"

	"backend/internal/config"
	"backend/internal/repository/postgres"
	"backend/internal/service/engine"
	"backend/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type App struct {
	router       *chi.Mux
	cfg          config.Config
	db           *sqlx.DB
	scenarioRepo *postgres.ScenarioRepo
	engine       *engine.Engine
	stateRepo    *postgres.StateRepo
}

func New() (*App, error) {
	cfg := config.Load()
	db := database.NewPostgres(cfg.PostgresDSN)
	graphRepo := postgres.NewGraphRepo(db)
	stateRepo := postgres.NewStateRepo(db)
	nodes, _ := graphRepo.GetNodes()
	edges, _ := graphRepo.GetEdges()

	eng := engine.New(nodes, edges)

	scenarioRepo := postgres.NewScenarioRepo(db)

	r := chi.NewRouter()

	app := &App{
		router:       r,
		cfg:          cfg,
		db:           db,
		scenarioRepo: scenarioRepo,
		engine:       eng,
		stateRepo:    stateRepo,
	}

	app.initRoutes()

	return app, nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":"+a.cfg.Port, a.router)
}
