package handlers

import (
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"io"
	"net/http"
	"net/url"
)

func (a *App) CompressHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &a.Config.CookieKey
	if *cookie == "" {
		fmt.Println("cookie is empty")
		auth.CreateNewUser(w)
		//w.WriteHeader(http.StatusUnauthorized)
		//return
	}

	fmt.Println("CompressHandler COOKIE = ", cookie)

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		logger.Errorf("body is nil or empty: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.Config.DatabaseDSN != "" {
		short, err := a.Storage.CheckIsURLExists(string(body))
		if err != nil {
			logger.Errorf("error is CheckIsURLExists: %s", err)
		}

		if short != "" {
			res, err := url.JoinPath(a.Config.BaseShortURL, short)
			if err != nil {
				logger.Errorf("error is JoinPath: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(res))
			if err != nil {
				logger.Errorf("Failed to send URL: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			return
		}
	}

	short, err := a.Storage.Save(string(body), "", a.Config.CookieKey)
	if err != nil {
		logger.Errorf("storage save is error: %s", err)
		w.WriteHeader(http.StatusBadRequest) //TODO в будущем переделать на http.StatusInternalServerError
		return
	}

	//fmt.Println("UserURLsHandler COOKIE = ", *cookie)
	//l, s := a.Storage.GetAllURLs(*cookie)

	//fmt.Println("UserURLsHandler long = ", l)
	//fmt.Println("UserURLsHandler short = ", s)

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
