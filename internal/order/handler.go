package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
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

func (h *Handler) Checkout(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := auth.GetUserFromContext(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	order, err := h.repository.Checkout(
		user.ID,
		user.Name,
		user.Email,
	)
	if err != nil {

		// Tampilkan error asli di terminal
		log.Printf(
			"CHECKOUT ERROR: %v",
			err,
		)

		if errors.Is(
			err,
			ErrCartEmpty,
		) {
			http.Error(
				w,
				"cart is empty",
				http.StatusBadRequest,
			)
			return
		}

		http.Error(
			w,
			"failed to checkout",
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

	json.NewEncoder(w).Encode(
		order,
	)
}

func (h *Handler) HandleNotification(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload struct {
		OrderID           string `json:"order_id"`
		TransactionID     string `json:"transaction_id"`
		TransactionStatus string `json:"transaction_status"`
		PaymentType       string `json:"payment_type"`
		FraudStatus       string `json:"fraud_status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	orderIDValue := strings.TrimPrefix(payload.OrderID, "bukuin-")
	if idx := strings.Index(orderIDValue, "-"); idx >= 0 {
		orderIDValue = orderIDValue[:idx]
	}
	orderID, err := strconv.Atoi(orderIDValue)
	if err != nil || orderID <= 0 {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	status := strings.ToLower(payload.TransactionStatus)
	if status == "capture" || status == "settlement" {
		status = "paid"
	}
	if status == "pending" {
		status = "pending"
	}
	if status == "cancel" || status == "expire" || status == "deny" {
		status = "cancelled"
	}
	if status == "refund" || status == "partial_refund" {
		status = "refunded"
	}

	paidAt := time.Now()
	if status != "paid" {
		paidAt = time.Time{}
	}

	if err := h.repository.UpdatePaymentStatus(
		orderID,
		status,
		payload.PaymentType,
		payload.TransactionID,
		paidAt,
	); err != nil {
		log.Printf("UPDATE PAYMENT STATUS ERROR: %v", err)
		http.Error(w, "failed to update payment status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetOrders(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := auth.GetUserFromContext(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	orders, err :=
		h.repository.GetOrdersByUserID(
			user.ID,
		)

	if err != nil {
		log.Printf(
			"GET ORDERS ERROR: %v",
			err,
		)

		http.Error(
			w,
			"failed to get orders",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		orders,
	)
}

func (h *Handler) GetOrderByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := auth.GetUserFromContext(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	orderID, err := strconv.Atoi(
		chi.URLParam(
			r,
			"id",
		),
	)
	if err != nil {
		http.Error(
			w,
			"invalid order id",
			http.StatusBadRequest,
		)
		return
	}

	order, err :=
		h.repository.GetOrderByID(
			user.ID,
			orderID,
		)

	if err != nil {

		log.Printf(
			"GET ORDER ERROR: %v",
			err,
		)

		if errors.Is(
			err,
			sql.ErrNoRows,
		) {
			http.Error(
				w,
				"order not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"failed to get order",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		order,
	)
}
