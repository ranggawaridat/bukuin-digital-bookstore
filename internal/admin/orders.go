package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
)

type OrderSummary struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	UserName    string  `json:"user_name"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

func (h *Handler) Orders(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	rows, err := h.db.Query(`
		SELECT
			o.id,
			o.user_id,
			u.name,
			o.total_amount,
			o.status
		FROM orders o
		JOIN users u ON u.id = o.user_id
		ORDER BY o.created_at DESC
	`)
	if err != nil {
		http.Error(w, "failed to get orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []OrderSummary{}
	for rows.Next() {
		var order OrderSummary
		if err := rows.Scan(&order.ID, &order.UserID, &order.UserName, &order.TotalAmount, &order.Status); err != nil {
			http.Error(w, "failed to get orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) OrderDetails(w http.ResponseWriter, r *http.Request) {
	_ = sql.ErrNoRows
	_ = json.NewEncoder
	_ = http.StatusOK
}
