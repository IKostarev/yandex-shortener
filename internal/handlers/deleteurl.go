package handlers

import (
	"encoding/json"
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

	if len(urls) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type DeleteResult struct {
		ShortURL string
		Success  bool
	}

	deleteURL := func(shortURL string, results chan<- DeleteResult) {
		err := a.Storage.DeleteURL([]byte(shortURL), string(cookie))
		results <- DeleteResult{ShortURL: shortURL, Success: err == true} //TODO maybe need err == false
	}

	results := make(chan DeleteResult)
	done := make(chan struct{})

	for _, shortURL := range urls {
		go deleteURL(shortURL, results)
	}

	fanIn := func(results <-chan DeleteResult, done chan<- struct{}) <-chan DeleteResult {
		merged := make(chan DeleteResult)

		go func() {
			defer close(merged)
			defer close(done)

			for result := range results {
				merged <- result
			}
		}()

		return merged
	}

	mergedResults := fanIn(results, done)

	success := true
	for range urls {
		result := <-mergedResults
		if !result.Success {
			success = false
		}
	}

	if success {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
