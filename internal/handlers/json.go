package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/authorization"
	"github.com/google/uuid"
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

	fmt.Println("61 СТРОКА Я ЗДЕСЬ")

	value := r.Context().Value(authorization.ContextKey("userID")).(string)
	fmt.Println("JSON value = ", value)
	user, err := uuid.Parse(value)
	fmt.Println("JSON user = ", user)
	if err != nil {
		logger.Errorf("error parse user uuid is: %s", err)
	}

	fmt.Println("72 СТРОКА Я ЗДЕСЬ")

	short, err := a.Storage.Save(req.ServerURL, "", user)
	if err != nil {
		logger.Errorf("storage save is error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("81 СТРОКА Я ЗДЕСЬ")

	resp.BaseShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
	if err != nil {
		logger.Errorf("join path have err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("90 СТРОКА Я ЗДЕСЬ")

	respContent, err := json.Marshal(resp)
	if err != nil {
		logger.Errorf("json marshal is error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("99 СТРОКА Я ЗДЕСЬ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(respContent); err != nil {
		logger.Errorf("Failed to send URL on json handler: %s", err)
	}
}
