package handlers

import (
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/middleware/auth"
	"net/http"
)

func (a *App) UserURLsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &a.Config.CookieKey
	if *cookie == "" {
		auth.CreateNewUser(w)
	}

	fmt.Println("UserURLsHandler COOKIE = ", *cookie)
	long, short := a.Storage.GetAllURLs(*cookie)

	fmt.Println("UserURLsHandler long = ", long)
	fmt.Println("UserURLsHandler short = ", short)
}
