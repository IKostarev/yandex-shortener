package handlers

import (
	"encoding/json"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
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
		logger.Errorf("json decode error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.Config.DatabaseDSN != "" {
		if err := a.checkURLExistsAndRespondJSON(w, req.ServerURL, &resp); err != nil {
			logger.Errorf("error checking URL existence: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	short, err := a.Storage.Save(req.ServerURL, "")
	if err != nil {
		logger.Errorf("storage save error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.BaseShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
	if err != nil {
		logger.Errorf("join path error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respContent, err := json.Marshal(resp)
	if err != nil {
		logger.Errorf("json marshal error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(respContent); err != nil {
		logger.Errorf("Failed to send URL on json handler: %s", err)
	}
}

func (a *App) checkURLExistsAndRespondJSON(w http.ResponseWriter, serverURL string, resp *ResultResponse) error {
	short, err := a.Storage.CheckIsURLExists(serverURL)
	if err != nil {
		logger.Errorf("error CheckIsURLExists on checkURLExistsAndRespondJSON: %s", err)
		return err
	}

	if short != "" {
		resp.BaseShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
		if err != nil {
			logger.Errorf("error JoinPath on checkURLExistsAndRespondJSON: %s", err)
			return err
		}

		respContent, err := json.Marshal(resp)
		if err != nil {
			logger.Errorf("error Marshal on checkURLExistsAndRespondJSON: %s", err)
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if _, err := w.Write(respContent); err != nil {
			logger.Errorf("error Write on checkURLExistsAndRespondJSON: %s", err)
			return err
		}

		return nil
	}

	return nil
}
