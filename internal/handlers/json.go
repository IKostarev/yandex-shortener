package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"net/http"
	"net/url"
)

type URLRequest struct {
	ServerURL string `json:"url"`
}

type ResultResponse struct {
	BaseShortURL string `json:"result"`
}

func (a *App) JSONHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("ID")
	if cookie == nil {
		cookie = auth.CreateNewUser(w)
		w.WriteHeader(http.StatusUnauthorized)
	}

	fmt.Println("JSONHandler COOKIE = ", cookie)

	var req URLRequest
	var resp ResultResponse

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("json decode is error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.Config.DatabaseDSN != "" {
		short, err := a.Storage.CheckIsURLExists(req.ServerURL)
		if err != nil {
			logger.Errorf("error is CheckIsURLExists: %s", err)
		}

		if short != "" {
			resp.BaseShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
			if err != nil {
				logger.Errorf("error is JoinPath: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			respContent, err := json.Marshal(resp)
			if err != nil {
				logger.Errorf("json marshal is error: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write(respContent); err != nil {
				logger.Errorf("Failed to send URL on json handler: %s", err)
			}
			return
		}
	}

	short, err := a.Storage.Save(req.ServerURL, "")
	if err != nil {
		logger.Errorf("storage save is error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.BaseShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
	if err != nil {
		logger.Errorf("join path have err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respContent, err := json.Marshal(resp)
	if err != nil {
		logger.Errorf("json marshal is error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(respContent); err != nil {
		logger.Errorf("Failed to send URL on json handler: %s", err)
	}
}
