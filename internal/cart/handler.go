package cart

import (
	"encoding/json"
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

func (h *Handler) GetCart(
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

	cart, err := h.repository.GetCart(
		user.ID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get cart",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(cart)
}

func (h *Handler) AddItem(
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

	var request AddItemRequest

	err = json.NewDecoder(
		r.Body,
	).Decode(&request)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	cart, err := h.repository.GetOrCreateCart(
		user.ID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get cart",
			http.StatusInternalServerError,
		)
		return
	}

	err = h.repository.AddItem(
		cart.ID,
		request.BookID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to add item",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "book added to cart",
		},
	)
}

func (h *Handler) UpdateItem(
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

	itemID, err := strconv.Atoi(
		chi.URLParam(r, "id"),
	)
	if err != nil {
		http.Error(
			w,
			"invalid item id",
			http.StatusBadRequest,
		)
		return
	}

	var request UpdateItemRequest

	err = json.NewDecoder(
		r.Body,
	).Decode(&request)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.Quantity < 1 {
		http.Error(
			w,
			"quantity must be at least 1",
			http.StatusBadRequest,
		)
		return
	}

	cart, err := h.repository.GetOrCreateCart(
		user.ID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get cart",
			http.StatusInternalServerError,
		)
		return
	}

	err = h.repository.UpdateItemQuantity(
		cart.ID,
		itemID,
		request.Quantity,
	)
	if err != nil {
		http.Error(
			w,
			"failed to update item",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "cart updated",
		},
	)
}

func (h *Handler) DeleteItem(
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

	itemID, err := strconv.Atoi(
		chi.URLParam(r, "id"),
	)
	if err != nil {
		http.Error(
			w,
			"invalid item id",
			http.StatusBadRequest,
		)
		return
	}

	cart, err := h.repository.GetOrCreateCart(
		user.ID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get cart",
			http.StatusInternalServerError,
		)
		return
	}

	err = h.repository.DeleteItem(
		cart.ID,
		itemID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to delete item",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "item deleted",
		},
	)
}
