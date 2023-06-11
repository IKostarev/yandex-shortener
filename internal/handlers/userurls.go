package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/authorization"
	"github.com/IKostarev/yandex-go-dev/internal/model"
	"github.com/google/uuid"
	"net/http"
)

func getUser(r *http.Request) (user uuid.UUID, err error) {
	userID, ok := r.Context().Value(authorization.ContextKey("userID")).(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("value is not string: %s", userID)
	}

	return uuid.Parse(userID)
}

func (a *App) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r) //TODO handle error

	links, _ := a.Storage.GetUserLinks(user) //TODO handle error

	resp := make([]model.UserLink, 0)
	for _, link := range links {
		resp = append(resp, model.UserLink{
			OriginalURL: link.OriginalURL,
			ShortURL:    fmt.Sprintf("%s/%s", a.Config.BaseShortURL, link.ShortURL),
		})
	}

	respJSON, _ := json.Marshal(resp) //TODO handle error

	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write(respJSON) //TODO handle error
	}
}
