package authorization

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ContextKey string

func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("USER") //TODO handle error
		if err != nil {
			fmt.Println("COOKIE IS ERROR = ", cookie)
			//w.WriteHeader(http.StatusUnauthorized)
			setNewUser(next, w, r, createNewUser(w))
			return
		}

		payload, err := base64.StdEncoding.DecodeString(cookie.Value) //TODO handle error
		if err != nil {
			fmt.Println("PAYLOAD IS ERROR = ", cookie)
			//setNewUser(next, w, r, createNewUser(w))
			return
		}

		h := hmac.New(sha256.New, []byte("empty"))
		h.Write(payload[:16])

		signed := h.Sum(nil)

		if !hmac.Equal(signed, payload[16:]) {
			fmt.Println("EQUAL IS ERROR = ", cookie)
			//setNewUser(next, w, r, createNewUser(w))
		}

		user, _ := uuid.FromBytes(payload[:16]) //TODO handle error

		fmt.Println("дошел до конца IS ERROR = ", cookie)
		setNewUser(next, w, r, user.String())
	})
}

func setNewUser(next http.Handler, w http.ResponseWriter, r *http.Request, user string) {
	ctx := context.WithValue(r.Context(), ContextKey("userID"), user)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func createNewUser(w http.ResponseWriter) string {
	user := uuid.New()

	b, _ := user.MarshalBinary()

	h := hmac.New(sha256.New, []byte("empty"))
	h.Write(b)

	signed := h.Sum(nil)

	cookie := http.Cookie{
		Name:    "USER",
		Value:   base64.StdEncoding.EncodeToString(append(b, signed...)),
		Expires: time.Now().Add(356 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)

	return user.String()
}
