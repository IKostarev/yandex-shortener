package auth

import (
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/config"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("ID")
		fmt.Println("IT'S COOKIE = ", cookie.Value)
		if cookie.Value == "" {
			CreateNewUser(w)
		}

		next.ServeHTTP(w, r)
	})
}

func CreateNewUser(w http.ResponseWriter) *http.Cookie {
	user := uuid.New().String()
	cfg := config.Config{}

	newCookie := http.Cookie{
		Name:    "ID",
		Value:   user,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	fmt.Println("CREATE NEW USER = ", user)

	fmt.Println("&newCookie = ", &newCookie)

	http.SetCookie(w, &newCookie)
	cfg.CookieKey = user

	fmt.Println("CREATE NEW USER COOKIE = ", cfg.CookieKey)

	return &newCookie
}
