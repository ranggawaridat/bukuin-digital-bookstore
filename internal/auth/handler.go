package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	repository *Repository
}

func NewHandler(
	repository *Repository,
) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request RegisterRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	request.Name = strings.TrimSpace(
		request.Name,
	)

	request.Email = strings.TrimSpace(
		strings.ToLower(request.Email),
	)

	if request.Name == "" {
		http.Error(
			w,
			"name is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Email == "" {
		http.Error(
			w,
			"email is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Password == "" {
		http.Error(
			w,
			"password is required",
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.repository.CreateUser(
		request.Name,
		request.Email,
		request.Password,
	)

	if err != nil {
		http.Error(
			w,
			"failed to create user",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request LoginRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	request.Email = strings.TrimSpace(
		strings.ToLower(request.Email),
	)

	user, err := h.repository.GetUserByEmail(
		request.Email,
	)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"failed to login",
			http.StatusInternalServerError,
		)
		return
	}

	if user.Password != request.Password {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)
		return
	}

	token, err := GenerateSessionToken()
	if err != nil {
		http.Error(
			w,
			"failed to create session",
			http.StatusInternalServerError,
		)
		return
	}

	expiresAt := time.Now().Add(
		24 * time.Hour,
	)

	err = h.repository.CreateSession(
		user.ID,
		token,
		expiresAt.UTC().Format(
			"2006-01-02 15:04:05",
		),
	)

	if err != nil {
		http.Error(
			w,
			"failed to create session",
			http.StatusInternalServerError,
		)
		return
	}

	SetSessionCookie(
		w,
		token,
		expiresAt,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "login successful",
		},
	)
}

func (h *Handler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	cookie, err := r.Cookie(
		SessionCookieName,
	)

	if err == nil {
		err = h.repository.DeleteSession(
			cookie.Value,
		)

		if err != nil {
			http.Error(
				w,
				"failed to logout",
				http.StatusInternalServerError,
			)
			return
		}
	}

	ClearSessionCookie(w)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "logout successful",
		},
	)
}

func (h *Handler) Me(
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

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"failed to get user",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(user)
}
