package auth

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

func (h *Handler) RequireAuth(
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		cookie, err := r.Cookie(
			SessionCookieName,
		)

		if err != nil {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		user, err := h.repository.GetUserBySession(
			cookie.Value,
		)

		if err != nil {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			user,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	}
}

func GetUserFromContext(
	ctx context.Context,
) (*User, error) {
	user, ok := ctx.Value(
		UserContextKey,
	).(*User)

	if !ok {
		return nil, errors.New(
			"user not found in context",
		)
	}

	return user, nil
}
