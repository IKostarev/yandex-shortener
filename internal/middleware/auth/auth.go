package auth

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("ID")
		if cookie == nil {
			CreateNewUser(w)
		}

		next.ServeHTTP(w, r)
	})
}

func CreateNewUser(w http.ResponseWriter) *http.Cookie {
	user := uuid.New().String()

	newCookie := http.Cookie{
		Name:    "ID",
		Value:   user,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	http.SetCookie(w, &newCookie)

	return &newCookie
}
