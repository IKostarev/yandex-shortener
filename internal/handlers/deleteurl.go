package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"io"
	"net/http"
)

func (a *App) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := auth.GlobalCookieKey

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var urls []string
	err = json.Unmarshal(body, &urls)
	if err != nil {
		logger.Errorf("failed to unmarshal request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("URLS = ", urls)

	if len(urls) == 0 {
		//logger.Error("empty list of URLs received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	success := true
	//ctx := r.Context()
	for _, shortURL := range urls {
		err := a.Storage.DeleteURL([]byte(shortURL), string(cookie))
		if !err {
			//logger.Errorf("failed to delete URL '%s': %s", shortURL, err)
			success = false
		}
	}

	if success {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
