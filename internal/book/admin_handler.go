package book

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	request, err := parseCreateBookRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	request, err := parseCreateBookRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func parseCreateBookRequest(r *http.Request) (CreateBookRequest, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return CreateBookRequest{}, err
	}

	request := CreateBookRequest{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Author:      strings.TrimSpace(r.FormValue("author")),
		Category:    strings.TrimSpace(r.FormValue("category")),
		Description: strings.TrimSpace(r.FormValue("description")),
		CoverURL:    strings.TrimSpace(r.FormValue("cover_url")),
		FilePath:    strings.TrimSpace(r.FormValue("file_path")),
	}

	price, err := strconv.Atoi(r.FormValue("price"))
	if err == nil {
		request.Price = price
	}

	cover, header, err := r.FormFile("cover")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return request, nil
		}
		return request, err
	}
	defer cover.Close()

	contentType := strings.TrimSpace(header.Header.Get("Content-Type"))
	if contentType != "" && !strings.HasPrefix(contentType, "image/") {
		return request, errors.New("cover must be an image file")
	}

	uploadDir := filepath.Clean("./web/static/uploads")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return request, err
	}

	fileExt := filepath.Ext(header.Filename)
	if fileExt == "" {
		fileExt = ".png"
	}

	fileName := fmt.Sprintf("book-cover-%d%s", time.Now().UnixNano(), fileExt)
	destination := filepath.Join(uploadDir, fileName)

	destinationFile, err := os.Create(destination)
	if err != nil {
		return request, err
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, cover); err != nil {
		return request, err
	}

	request.CoverURL = "/static/uploads/" + fileName
	return request, nil
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
