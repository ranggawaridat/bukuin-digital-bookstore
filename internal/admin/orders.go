package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
)

type OrderSummary struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	UserName      string     `json:"user_name"`
	TotalAmount   float64    `json:"total_amount"`
	Status        string     `json:"status"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	TransactionID string     `json:"transaction_id,omitempty"`
	PaymentURL    string     `json:"payment_url,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
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
			o.status,
			o.payment_method,
			o.transaction_id,
			o.payment_url,
			o.paid_at
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
		var paymentMethod sql.NullString
		var transactionID sql.NullString
		var paymentURL sql.NullString
		var paidAt sql.NullTime

		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.UserName,
			&order.TotalAmount,
			&order.Status,
			&paymentMethod,
			&transactionID,
			&paymentURL,
			&paidAt,
		); err != nil {
			http.Error(w, "failed to get orders", http.StatusInternalServerError)
			return
		}

		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		if transactionID.Valid {
			order.TransactionID = transactionID.String
		}
		if paymentURL.Valid {
			order.PaymentURL = paymentURL.String
		}
		if paidAt.Valid {
			paidAtCopy := paidAt.Time
			order.PaidAt = &paidAtCopy
		}

		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || orderID <= 0 {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	status := strings.ToLower(strings.TrimSpace(payload.Status))
	if status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	allowedStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"cancelled": true,
		"refunded":  true,
	}

	if !allowedStatuses[status] {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	var paidAtValue any
	if status == "paid" {
		paidAtValue = time.Now()
	}

	result, err := h.db.Exec(`
		UPDATE orders
		SET status = ?, paid_at = COALESCE(paid_at, ?)
		WHERE id = ?
	`, status, paidAtValue, orderID)
	if err != nil {
		http.Error(w, "failed to update order status", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to update order status", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
