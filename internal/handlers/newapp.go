package handlers

import (
	"github.com/IKostarev/yandex-go-dev/internal/config"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/gzip"
	"github.com/IKostarev/yandex-go-dev/internal/storage"
	"github.com/go-chi/chi/v5"
)

type App struct {
	*chi.Mux
	Config  config.Config
	Storage storage.Storage
}

func NewApp(cfg config.Config, store storage.Storage) *App {
	app := App{
		chi.NewRouter(),
		cfg,
		store,
	}

	app.Use(gzip.Request)
	app.Use(gzip.Response)

	app.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(app.CompressHandler))
		r.Get("/{id}", logger.ResponseLogger(app.GetURLHandler))
		r.Get("/ping", logger.ResponseLogger(app.PingHandler))

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", app.JSONHandler)
			r.Post("/shorten/batch", app.BatchHandler)
		})
	})

	return &app
}
