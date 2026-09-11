package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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
