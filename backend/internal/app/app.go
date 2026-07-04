package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	router *chi.Mux
}

func New() (*App, error) {
	r := chi.NewRouter()

	app := &App{
		router: r,
	}

	app.initRoutes()

	return app, nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":8080", a.router)
}
