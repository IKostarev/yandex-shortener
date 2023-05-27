package handlers

import (
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"io"
	"net/http"
	"net/url"
)

func (a *App) CompressHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		logger.Errorf("body is nil or empty: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.Config.DatabaseDSN != "" {
		if err := a.checkURLExistsAndRespond(w, body); err != nil {
			logger.Errorf("error checking URL existence: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	short, err := a.Storage.Save(string(body), "")
	if err != nil {
		logger.Errorf("storage save error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	long, err := url.JoinPath(a.Config.BaseShortURL, short)
	if err != nil {
		logger.Errorf("join path error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(long))
	if err != nil {
		logger.Errorf("Failed to send URL: %s", err)
	}
}

func (a *App) checkURLExistsAndRespond(w http.ResponseWriter, body []byte) error {
	short, err := a.Storage.CheckIsURLExists(string(body))
	if err != nil {
		logger.Errorf("error CheckIsURLExists on checkURLExistsAndRespond: %s", err)
		return err
	}

	if short != "" {
		res, err := url.JoinPath(a.Config.BaseShortURL, short)
		if err != nil {
			logger.Errorf("error JoinPath on checkURLExistsAndRespond: %s", err)
			return err
		}

		_, err = w.Write([]byte(res))
		if err != nil {
			logger.Errorf("error Write on checkURLExistsAndRespond: %s", err)
			return err
		}

		w.WriteHeader(http.StatusConflict)
		return nil
	}

	return nil
}
