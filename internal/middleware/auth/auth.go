package auth

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("ID")
		if cookie == nil {
			fmt.Println("Cookie = ", cookie)

			fmt.Println("createNewUser(w) = ", createNewUser(w))

			//w.WriteHeader(http.StatusUnauthorized)
			//return
		}

		fmt.Println("Cookie NOT NIL = ", cookie)
		next.ServeHTTP(w, r)
	})
}

func createNewUser(w http.ResponseWriter) string {
	user := uuid.New().String()

	newCookie := http.Cookie{
		Name:    "ID",
		Value:   user,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	http.SetCookie(w, &newCookie)

	return user
}
