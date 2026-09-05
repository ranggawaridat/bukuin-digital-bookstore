package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

const SessionCookieName = "bukuin_session"

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func SetSessionCookie(
	w http.ResponseWriter,
	token string,
	expiresAt time.Time,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     SessionCookieName,
			Value:    token,
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
}

func ClearSessionCookie(
	w http.ResponseWriter,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     SessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
}
