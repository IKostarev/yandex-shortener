package auth

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CookieKey string

var GlobalCookieKey CookieKey

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

func CreateNewUser(_ http.ResponseWriter) *http.Cookie {
	user := uuid.New()

	newCookie := http.Cookie{
		Name:    "ID",
		Value:   user.String(),
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	GlobalCookieKey = CookieKey(newCookie.Value)

	return &newCookie
}
