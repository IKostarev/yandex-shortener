package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"log"
	"net/http"
	"strconv"
)

type UserLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (a *App) UserURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := a.Config.CookieKey
	if cookie == "" {
		auth.CreateNewUser(w)
		//w.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Println("UserURLsHandler = ", cookie)

	links, _ := a.Storage.GetAllURLs(cookie)

	response := make([]UserLink, 0)
	for _, link := range links {
		response = append(response, UserLink{
			ShortURL:    fmt.Sprintf("%s/%d", a.Config.BaseShortURL, link[0]),
			OriginalURL: strconv.Itoa(int(link[1])),
		})
	}

	_, err := json.Marshal(response)
	if err != nil {
		//log.Printf("unable to marshal response: %v", err)
		return
	}

	if len(response) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set("content-type", "application/json")
		//w.Header().Set("X-Content-Type-Options", "nosniff")
		//_, err = w.Write(responseJSON)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
