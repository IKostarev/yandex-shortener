package handlers

import (
	"errors"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/authorization"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "id")
	if url == "" {
		_ = errors.New("url param bad with id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := r.Context().Value(authorization.ContextKey("userID")).(string)
	fmt.Println("GETURL value = ", value)
	user, err := uuid.Parse(value)
	fmt.Println("GETURL user = ", user)
	if err != nil {
		logger.Errorf("error parse user uuid is: %s", err)
	}

	m, _ := a.Storage.Get(url, "", user)
	if m == "" {
		logger.Errorf("get url is bad: %s", url)
		w.WriteHeader(http.StatusBadRequest) //TODO в будущем переделать на http.StatusNotFound
		return
	}

	w.Header().Add("Location", m)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
