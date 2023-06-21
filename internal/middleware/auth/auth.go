package auth

import (
	"fmt"
	"net/http"
)

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("ID")
		fmt.Println("Cookie = ", cookie)
		next.ServeHTTP(w, r)
	})
}
