package handlers

import (
	"errors"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "id")
	if url == "" {
		_ = errors.New("url param bad with id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, _ := a.Storage.Get(url, "")
	if m == "" {
		logger.Errorf("get url is bad: %s", url)
		w.WriteHeader(http.StatusBadRequest) //TODO в будущем переделать на http.StatusNotFound
		return
	}

	w.Header().Add("Location", m)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
