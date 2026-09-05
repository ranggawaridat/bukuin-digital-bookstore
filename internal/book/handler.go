package book

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) GetBooks(
	w http.ResponseWriter,
	r *http.Request,
) {
	search := r.URL.Query().Get("q")

	category := r.URL.Query().Get("category")

	books, err := h.repository.GetAll(
		search,
		category,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get books",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(books)
}

func (h *Handler) GetBookByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := chi.URLParam(
		r,
		"id",
	)

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(
			w,
			"invalid book id",
			http.StatusBadRequest,
		)
		return
	}

	book, err := h.repository.GetByID(id)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(
			w,
			"book not found",
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"failed to get book",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(book)
}
