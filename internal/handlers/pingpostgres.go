package handlers

import (
	"fmt"
	"net/http"
)

func (a *App) PingHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &a.Config.CookieKey

	fmt.Println("PingHandler COOKIE = ", cookie)

	if !a.Storage.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
