package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"io"
	"net/http"
	"time"
)

func (a *App) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := auth.GlobalCookieKey

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("ошибка чтения тела запроса: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var urls []string
	err = json.Unmarshal(body, &urls)
	if err != nil {
		logger.Errorf("ошибка разбора тела запроса: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	go func() {
		time.Sleep(3 * time.Minute)

		success := a.Storage.DeleteURL(urls, string(cookie))

		if success {
			fmt.Println("URL-ы успешно удалены")
		} else {
			logger.Errorf("ошибка при удалении URL-ов")
		}
	}()
}
