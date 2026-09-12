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

	existingBook, err := h.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to load book", http.StatusInternalServerError)
		return
	}

	if request.CoverURL == "" {
		request.CoverURL = existingBook.CoverURL
	}
	if request.FilePath == "" {
		request.FilePath = existingBook.FilePath
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
	if err := r.ParseMultipartForm(50 << 20); err != nil {
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

	uploadDir := filepath.Clean("./web/static/uploads")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return request, err
	}

	cover, coverHeader, err := r.FormFile("cover")
	if err == nil {
		defer cover.Close()

		contentType := strings.TrimSpace(coverHeader.Header.Get("Content-Type"))
		if contentType != "" && !strings.HasPrefix(contentType, "image/") {
			return request, errors.New("cover must be an image file")
		}

		fileExt := filepath.Ext(coverHeader.Filename)
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
	} else if !errors.Is(err, http.ErrMissingFile) {
		return request, err
	}

	ebook, ebookHeader, err := r.FormFile("ebook")
	if err == nil {
		defer ebook.Close()

		contentType := strings.TrimSpace(ebookHeader.Header.Get("Content-Type"))
		isPDFContentType := strings.Contains(strings.ToLower(contentType), "pdf") || strings.HasSuffix(strings.ToLower(ebookHeader.Filename), ".pdf")
		if contentType != "" && !isPDFContentType {
			return request, errors.New("ebook must be a PDF file")
		}

		fileExt := filepath.Ext(ebookHeader.Filename)
		if fileExt == "" || !strings.EqualFold(fileExt, ".pdf") {
			fileExt = ".pdf"
		}

		fileName := fmt.Sprintf("book-file-%d%s", time.Now().UnixNano(), fileExt)
		destination := filepath.Join(uploadDir, fileName)

		destinationFile, err := os.Create(destination)
		if err != nil {
			return request, err
		}
		defer destinationFile.Close()

		if _, err := io.Copy(destinationFile, ebook); err != nil {
			return request, err
		}

		request.FilePath = "/static/uploads/" + fileName
	} else if !errors.Is(err, http.ErrMissingFile) {
		return request, err
	}

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
