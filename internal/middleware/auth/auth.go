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
		//fmt.Println("IT'S COOKIE = ", cookie.Value)
		if cookie == nil {
			_ = CreateNewUser(w)
		}

		next.ServeHTTP(w, r)
	})
}

func CreateNewUser(w http.ResponseWriter) *http.Cookie {
	user := uuid.New()
	cfg := config.Config{}

	newCookie := http.Cookie{
		Name:    "ID",
		Value:   user.String(),
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	//fmt.Println("CREATE NEW USER = ", user)
	//
	//fmt.Println("CREATE NEW &USER = ", &user)
	//
	//fmt.Println("&newCookie = ", &newCookie)

	http.SetCookie(w, &newCookie)

	cfg.CookieKey = &newCookie

	//fmt.Println("CREATE NEW USER COOKIE = ", cfg.CookieKey)

	return &newCookie
}
