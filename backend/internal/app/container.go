package app

import (
	"net/http"
)

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	return srv.ListenAndServe()
}
