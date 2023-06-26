package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"net/http"
	"sync"
)

type UserLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (a *App) UserURLsHandler(w http.ResponseWriter, _ *http.Request) {
	cookie := auth.GlobalCookieKey

	var wg sync.WaitGroup

	var longURLs []string
	var shortURLs []string
	var longURLsErr, shortURLsErr string

	wg.Add(2)

	go func() {
		defer wg.Done()
		longURLs, longURLsErr = a.Storage.GetAllURLs(string(cookie))
	}()

	go func() {
		defer wg.Done()
		shortURLs, shortURLsErr = a.Storage.GetAllShortURLs(string(cookie))
	}()

	wg.Wait()

	if longURLsErr != "" || shortURLsErr != "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("longURLs = ", longURLs)
	fmt.Println("shortURLs = ", shortURLs)

	response := make([]UserLink, 0)
	for i := 0; i < len(shortURLs); i++ {
		shortURL := fmt.Sprintf("%s/%s", a.Config.BaseShortURL, shortURLs[i])
		originalURL := longURLs[i]

		response = append(response, UserLink{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(response) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseJSON)
	}
}
