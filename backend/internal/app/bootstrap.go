package app

import (
	"life-assistant/internal/config"
	"life-assistant/internal/graph"
	"life-assistant/internal/handler/http"
	"life-assistant/internal/repository/postgres"
	"life-assistant/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type App struct {
	router *chi.Mux
	db     *sqlx.DB
}

func New() (*App, error) {
	cfg := config.MustLoad()

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	repos := postgres.NewRepositories(db)

	graphEngine := graph.NewEngine(repos.Edge, repos.Node)

	services := service.NewServices(repos, graphEngine)

	router := http.NewRouter(services)

	return &App{
		router: router,
		db:     db,
	}, nil
}
