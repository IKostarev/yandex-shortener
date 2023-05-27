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

	//if a.Config.DatabaseDSN != "" { // проверка нужна чтоб отловить когда запущено приложение с бд и чтоб в бд был поиск уже сохраненного body
	//	check, err := a.Storage.CheckIsURLExists(string(body))
	//	if err != nil {
	//		logger.Errorf("Failed to CheckIsURLExists URL: %s", err)
	//		w.WriteHeader(http.StatusBadRequest)
	//	}
	//
	//	if check != "" {
	//		res, _ := url.JoinPath(a.Config.BaseShortURL, check)
	//
	//		w.WriteHeader(http.StatusConflict)
	//		_, err = w.Write([]byte(res))
	//		if err != nil {
	//			logger.Errorf("Failed to send URL: %s", err)
	//			w.WriteHeader(http.StatusBadRequest)
	//		}
	//		return
	//	}
	//}

	if a.Config.DatabaseDSN != "" {
		short, _ := a.Storage.CheckIsURLExists(string(body)) //TODO handle error
		if short != "" {
			res, _ := url.JoinPath(a.Config.BaseShortURL, short) //TODO handle error

			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(res))
			if err != nil {
				logger.Errorf("Failed to send URL: %s", err)
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
	}

	short, err := a.Storage.Save(string(body), "")
	if err != nil {
		logger.Errorf("storage save is error: %s", err)
		w.WriteHeader(http.StatusBadRequest) //TODO в будущем переделать на http.StatusInternalServerError
		return
	}

	long, err := url.JoinPath(a.Config.BaseShortURL, short)
	if err != nil {
		logger.Errorf("join path have err: %s", err)
		w.WriteHeader(http.StatusBadRequest) //TODO в будущем переделать на http.StatusInternalServerError
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(long))
	if err != nil {
		logger.Errorf("Failed to send URL: %s", err)
	}
}
