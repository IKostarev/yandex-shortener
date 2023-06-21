package handlers

import (
	"encoding/json"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/google/uuid"
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

func (a *App) BatchHandler(w http.ResponseWriter, rq *http.Request) {
	var req []URLsRequest
	var resp []URLsResponse

	if err := json.NewDecoder(rq.Body).Decode(&req); err != nil {
		logger.Errorf("json decode is error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, item := range req {
		var r URLsResponse

		short, err := a.Storage.CheckIsURLExists(item.OriginalURL)
		if err != nil {
			logger.Errorf("error CheckIsURLExists on BatchHandler: %s", err)
		}

		if short != "" {
			r.CorrelationID = item.CorrelationID
			r.ShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
			if err != nil {
				logger.Errorf("error JoinPath on BatchHandler: %s", err)
				w.WriteHeader(http.StatusBadRequest) // TODO: в будущем переделать на http.StatusInternalServerError
				return
			}

			resp = append(resp, r)
			w.WriteHeader(http.StatusConflict)
		} else {
			//value := rq.Context().Value(authorization.ContextKey("userID")).(string)

			//user, err := uuid.Parse(value)
			//if err != nil {
			//	logger.Errorf("error parse user uuid is: %s", err)
			//}
			user := uuid.New()

			short, err := a.Storage.Save(item.OriginalURL, item.CorrelationID, user)
			if err != nil {
				logger.Errorf("batch save is error: %s", err)
				w.WriteHeader(http.StatusBadRequest) // TODO: в будущем переделать на http.StatusInternalServerError
				return
			}

			r.CorrelationID = item.CorrelationID
			r.ShortURL, err = url.JoinPath(a.Config.BaseShortURL, short)
			if err != nil {
				logger.Errorf("join path has error: %s", err)
				w.WriteHeader(http.StatusBadRequest) // TODO: в будущем переделать на http.StatusInternalServerError
				return
			}

			resp = append(resp, r)
		}
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
		logger.Errorf("Failed to send URLsResponse on batch handler: %s", err)
	}
}
