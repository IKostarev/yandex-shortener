package handlers

import (
	"encoding/json"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"net/http"
	"net/url"
)

type URLsRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLsResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (a *App) BatchHandler(w http.ResponseWriter, r *http.Request) {
	var req []URLsRequest
	var resp []URLsResponse

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("json decode error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, item := range req {
		var r URLsResponse

		if err := a.checkURLExistsAndRespondBatch(item, &r); err != nil {
			logger.Errorf("error checking URL existence: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp = append(resp, r)
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
		logger.Errorf("Failed to send URLsResponse on batch handler: %s", err)
	}
}

func (a *App) checkURLExistsAndRespondBatch(req URLsRequest, resp *URLsResponse) error {
	short, err := a.Storage.CheckIsURLExists(req.OriginalURL)
	if err != nil {
		logger.Errorf("error CheckIsURLExists on checkURLExistsAndRespondBatch: %s", err)
		return err
	}

	if short != "" {
		resp.CorrelationID = req.CorrelationID
		resp.ShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
		if err != nil {
			logger.Errorf("error JoinPath on checkURLExistsAndRespondBatch: %s", err)
			return err
		}

		return nil
	}

	shortURL, err := a.Storage.Save(req.OriginalURL, req.CorrelationID)
	if err != nil {
		logger.Errorf("error Save on checkURLExistsAndRespondBatch: %s", err)
		return err
	}

	resp.CorrelationID = req.CorrelationID
	resp.ShortURL, err = url.JoinPath(a.Config.BaseShortURL, shortURL)
	if err != nil {
		logger.Errorf("error JoinPath on checkURLExistsAndRespondBatch: %s", err)
		return err
	}

	return nil
}
