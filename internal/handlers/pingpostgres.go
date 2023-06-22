package handlers

import (
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"net/http"
)

func (a *App) PingHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &a.Config.CookieKey
	if *cookie == "" {
		fmt.Println("cookie is empty")
		auth.CreateNewUser(w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("PingHandler COOKIE = ", cookie)

	if !a.Storage.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
