package handlers

import (
	"net/http"
)

func (a *App) PingHandler(w http.ResponseWriter, _ *http.Request) {
	if !a.Storage.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
