package handlers

import (
	"errors"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &a.Config.CookieKey
	if *cookie == "" {
		fmt.Println("cookie is empty")
		auth.CreateNewUser(w)
		//w.WriteHeader(http.StatusUnauthorized)
		//return
	}

	fmt.Println("GetURLHandler COOKIE = ", cookie)

	url := chi.URLParam(r, "id")
	if url == "" {
		_ = errors.New("url param bad with id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, _ := a.Storage.Get(url, "", a.Config.CookieKey)
	if m == "" {
		logger.Errorf("get url is bad: %s", url)
		w.WriteHeader(http.StatusBadRequest) //TODO в будущем переделать на http.StatusNotFound
		return
	}

	w.Header().Add("Location", m)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
