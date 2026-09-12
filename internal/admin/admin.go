package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
)

type Stats struct {
	TotalUsers   int     `json:"total_users"`
	TotalBooks   int     `json:"total_books"`
	TotalOrders  int     `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
}

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var stats Stats
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers); err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	if err := h.db.QueryRow(`SELECT COUNT(*) FROM books`).Scan(&stats.TotalBooks); err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	if err := h.db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&stats.TotalOrders); err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	if err := h.db.QueryRow(`SELECT COALESCE(SUM(total_amount), 0) FROM orders`).Scan(&stats.TotalRevenue); err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
