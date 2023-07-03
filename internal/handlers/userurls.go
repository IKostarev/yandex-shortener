package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"net/http"
)

type UserLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (a *App) UserURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := auth.GlobalCookieKey

	longURLs, _ := a.Storage.GetAllURLs(string(cookie))
	shortURLs, _ := a.Storage.GetAllShortURLs(string(cookie))

	fmt.Println("longURLs = ", longURLs)

	response := make([]UserLink, 0)
	for i := 0; i < len(longURLs) && i < len(shortURLs); i++ {
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
