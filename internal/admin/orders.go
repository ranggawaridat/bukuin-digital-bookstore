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
	UserEmail     string     `json:"user_email,omitempty"`
	TotalAmount   float64    `json:"total_amount"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
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
			u.email,
			o.total_amount,
			o.status,
			o.created_at,
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
			&order.UserEmail,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
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

func (h *Handler) getInvoice(orderID int) (map[string]any, error) {
	var order struct {
		ID            int
		UserID        int
		UserName      string
		UserEmail     string
		TotalAmount   float64
		Status        string
		CreatedAt     time.Time
		PaymentMethod sql.NullString
		TransactionID sql.NullString
		PaymentURL    sql.NullString
		PaidAt        sql.NullTime
	}

	if err := h.db.QueryRow(`
		SELECT
			o.id,
			o.user_id,
			u.name,
			u.email,
			o.total_amount,
			o.status,
			o.created_at,
			o.payment_method,
			o.transaction_id,
			o.payment_url,
			o.paid_at
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE o.id = ?
	`, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.UserName,
		&order.UserEmail,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.PaymentMethod,
		&order.TransactionID,
		&order.PaymentURL,
		&order.PaidAt,
	); err != nil {
		return nil, err
	}

	rows, err := h.db.Query(`
		SELECT
			oi.id,
			oi.book_id,
			oi.title,
			oi.author,
			oi.price,
			oi.quantity,
			oi.subtotal
		FROM order_items oi
		WHERE oi.order_id = ?
		ORDER BY oi.id ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []map[string]any{}
	for rows.Next() {
		var item struct {
			ID       int
			BookID   int
			Title    string
			Author   string
			Price    float64
			Quantity int
			Subtotal float64
		}
		if err := rows.Scan(
			&item.ID,
			&item.BookID,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.Quantity,
			&item.Subtotal,
		); err != nil {
			return nil, err
		}

		items = append(items, map[string]any{
			"id":       item.ID,
			"book_id":  item.BookID,
			"title":    item.Title,
			"author":   item.Author,
			"price":    item.Price,
			"quantity": item.Quantity,
			"subtotal": item.Subtotal,
		})
	}

	invoice := map[string]any{
		"id":             order.ID,
		"user_id":        order.UserID,
		"user_name":      order.UserName,
		"user_email":     order.UserEmail,
		"status":         order.Status,
		"total_amount":   order.TotalAmount,
		"created_at":     order.CreatedAt,
		"payment_method": order.PaymentMethod.String,
		"transaction_id": order.TransactionID.String,
		"payment_url":    order.PaymentURL.String,
		"paid_at":        nil,
		"items":          items,
	}

	if order.PaidAt.Valid {
		paidAtCopy := order.PaidAt.Time
		invoice["paid_at"] = paidAtCopy
	}

	return invoice, nil
}

func (h *Handler) Invoice(w http.ResponseWriter, r *http.Request) {
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

	invoice, err := h.getInvoice(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get invoice", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	stats := Stats{}
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

	rows, err := h.db.Query(`
		SELECT
			o.id,
			o.user_id,
			u.name,
			u.email,
			o.total_amount,
			o.status,
			o.created_at,
			o.payment_method,
			o.transaction_id,
			o.payment_url,
			o.paid_at
		FROM orders o
		JOIN users u ON u.id = o.user_id
		ORDER BY o.created_at DESC
	`)
	if err != nil {
		http.Error(w, "failed to get report", http.StatusInternalServerError)
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
			&order.UserEmail,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&paymentMethod,
			&transactionID,
			&paymentURL,
			&paidAt,
		); err != nil {
			http.Error(w, "failed to get report", http.StatusInternalServerError)
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

	response := map[string]any{
		"generated_at": time.Now(),
		"stats":        stats,
		"orders":       orders,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
