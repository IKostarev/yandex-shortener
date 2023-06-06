package authorization

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("USER") //TODO handle error

		payload, _ := base64.StdEncoding.DecodeString(cookie.Value) //TODO handle error

		h := hmac.New(sha256.New, []byte("empty"))
		h.Write(payload[:16])

		signed := h.Sum(nil)

		if !hmac.Equal(signed, payload[16:]) {
			setNewUser(next, w, r, createNewUser(w))
		}

		user, _ := uuid.FromBytes(payload[:16]) //TODO handle error

		setNewUser(next, w, r, user.String())
	})
}

func setNewUser(next http.Handler, w http.ResponseWriter, r *http.Request, user string) {
	ctx := context.WithValue(r.Context(), "userID", user)
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
