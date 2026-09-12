package book

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
)

func (h *Handler) CreateBook(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var request CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Title == "" || request.Author == "" || request.Category == "" || request.Description == "" {
		http.Error(w, "required fields are missing", http.StatusBadRequest)
		return
	}

	if request.Price < 0 {
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	book, err := h.repository.Create(request)
	if err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) UpdateBook(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	var request CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.repository.Update(id, request)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) DeleteBook(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	if err := h.repository.Delete(id); err != nil {
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "book deleted"})
}
