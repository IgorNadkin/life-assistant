package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *App) initRoutes() {
	a.router.Route("/api", func(r chi.Router) {

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		})
	})
}
