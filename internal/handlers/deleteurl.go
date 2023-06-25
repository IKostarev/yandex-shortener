package handlers

import (
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"io"
	"net/http"
)

func (a *App) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := auth.GlobalCookieKey

	urls, err := io.ReadAll(r.Body)
	if len(urls) == 0 || err != nil {
		logger.Errorf("body is nil or empty: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	del := a.Storage.DeleteURL(urls, string(cookie))
	if del == false {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}
