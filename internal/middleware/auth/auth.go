package auth

import (
	"github.com/IKostarev/yandex-go-dev/internal/config"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("ID")
		if cookie == nil {
			newCookie := CreateNewUser(w)
			http.SetCookie(w, newCookie)
		}

		next.ServeHTTP(w, r)
	})
}

func CreateNewUser(w http.ResponseWriter) *http.Cookie {
	user := uuid.New()
	cfg := &config.Config{}

	newCookie := http.Cookie{
		Name:    "ID",
		Value:   user.String(),
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	cfg.CookieKey = newCookie.Value

	return &newCookie
}
