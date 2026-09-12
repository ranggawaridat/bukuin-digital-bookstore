package order

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var ErrCartEmpty = errors.New(
	"cart is empty",
)

type Repository struct {
	db *sql.DB
}

func NewRepository(
	db *sql.DB,
) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Checkout(
	userID int,
	userName string,
	userEmail string,
) (*Order, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	committed := false

	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	var cartID int

	err = tx.QueryRow(
		`
		SELECT id
		FROM carts
		WHERE user_id = ?
		`,
		userID,
	).Scan(
		&cartID,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCartEmpty
	}

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(
		`
		SELECT
			cart_items.book_id,
			books.title,
			books.author,
			books.price,
			cart_items.quantity

		FROM cart_items

		JOIN books
			ON books.id = cart_items.book_id

		WHERE cart_items.cart_id = ?
		`,
		cartID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	type checkoutItem struct {
		BookID   int
		Title    string
		Author   string
		Price    float64
		Quantity int
		Subtotal float64
	}

	items := []checkoutItem{}

	var totalAmount float64

	for rows.Next() {

		var item checkoutItem

		err := rows.Scan(
			&item.BookID,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		item.Subtotal =
			item.Price *
				float64(item.Quantity)

		totalAmount +=
			item.Subtotal

		items = append(
			items,
			item,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, ErrCartEmpty
	}

	result, err := tx.Exec(
		`
		INSERT INTO orders (
			user_id,
			total_amount,
			status
		)
		VALUES (?, ?, ?)
		`,
		userID,
		totalAmount,
		"pending",
	)
	if err != nil {
		return nil, err
	}

	orderID, err :=
		result.LastInsertId()

	if err != nil {
		return nil, err
	}

	order := &Order{
		ID:          int(orderID),
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "pending",
		Items:       []OrderItem{},
	}

	for _, item := range items {

		result, err := tx.Exec(
			`
			INSERT INTO order_items (
				order_id,
				book_id,
				title,
				author,
				price,
				quantity,
				subtotal
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			`,
			orderID,
			item.BookID,
			item.Title,
			item.Author,
			item.Price,
			item.Quantity,
			item.Subtotal,
		)
		if err != nil {
			return nil, err
		}

		itemID, err :=
			result.LastInsertId()

		if err != nil {
			return nil, err
		}

		order.Items = append(
			order.Items,
			OrderItem{
				ID:       int(itemID),
				OrderID:  int(orderID),
				BookID:   item.BookID,
				Title:    item.Title,
				Author:   item.Author,
				Price:    item.Price,
				Quantity: item.Quantity,
				Subtotal: item.Subtotal,
			},
		)
	}

	_, err = tx.Exec(
		`
		DELETE FROM cart_items
		WHERE cart_id = ?
		`,
		cartID,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	committed = true

	if err := r.createMidtransPayment(
		order,
		userName,
		userEmail,
	); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *Repository) createMidtransPayment(
	order *Order,
	userName string,
	userEmail string,
) error {
	serverKey := strings.TrimSpace(
		os.Getenv("MIDTRANS_SERVER_KEY"),
	)
	if serverKey == "" {
		return nil
	}

	baseURL := strings.TrimSpace(
		os.Getenv("MIDTRANS_BASE_URL"),
	)
	if baseURL == "" {
		baseURL = "https://app.sandbox.midtrans.com"
	}

	itemDetails := make(
		[]map[string]any,
		0,
		len(order.Items),
	)
	midtransOrderID := fmt.Sprintf("bukuin-%d-%d", order.ID, time.Now().UnixNano())
	for _, item := range order.Items {
		itemDetails = append(
			itemDetails,
			map[string]any{
				"id":       fmt.Sprintf("%d", item.BookID),
				"price":    int(item.Price),
				"quantity": item.Quantity,
				"name":     item.Title,
			},
		)
	}

	payload := map[string]any{
		"transaction_details": map[string]any{
			"order_id":     midtransOrderID,
			"gross_amount": int(order.TotalAmount),
		},
		"item_details":     itemDetails,
		"enabled_payments": []string{"gopay", "bank_transfer", "shopeepay", "credit_card", "qris"},
	}

	if userName != "" || userEmail != "" {
		payload["customer_details"] = map[string]any{
			"first_name": userName,
			"email":      userEmail,
		}
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"POST",
		baseURL+"/snap/v1/transactions",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return err
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)
	request.Header.Set(
		"Accept",
		"application/json",
	)
	request.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(serverKey+":")),
	)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf(
			"midtrans request failed: %s",
			strings.TrimSpace(string(responseBody)),
		)
	}

	var snapResponse struct {
		Token         string `json:"token"`
		RedirectURL   string `json:"redirect_url"`
		TransactionID string `json:"transaction_id"`
	}

	if err := json.Unmarshal(
		responseBody,
		&snapResponse,
	); err != nil {
		return err
	}

	order.PaymentToken = snapResponse.Token
	order.PaymentURL = snapResponse.RedirectURL
	order.TransactionID = snapResponse.TransactionID

	return r.updateOrderPayment(
		order.ID,
		snapResponse.Token,
		snapResponse.RedirectURL,
		snapResponse.TransactionID,
	)
}

func (r *Repository) updateOrderPayment(
	orderID int,
	paymentToken string,
	paymentURL string,
	transactionID string,
) error {
	_, err := r.db.Exec(
		`
		UPDATE orders
		SET
			payment_token = ?,
			payment_url = ?,
			transaction_id = ?
		WHERE id = ?
		`,
		paymentToken,
		paymentURL,
		transactionID,
		orderID,
	)
	return err
}

func (r *Repository) UpdatePaymentStatus(
	orderID int,
	status string,
	paymentMethod string,
	transactionID string,
	paidAt time.Time,
) error {
	var paidAtValue any
	if !paidAt.IsZero() {
		paidAtValue = paidAt
	}

	_, err := r.db.Exec(
		`
		UPDATE orders
		SET
			status = ?,
			payment_method = ?,
			transaction_id = ?,
			paid_at = ?
		WHERE id = ?
		`,
		status,
		paymentMethod,
		transactionID,
		paidAtValue,
		orderID,
	)
	return err
}

func (r *Repository) GetOrdersByUserID(
	userID int,
) ([]Order, error) {

	rows, err := r.db.Query(
		`
		SELECT
			id,
			user_id,
			total_amount,
			status,
			created_at,
			payment_token,
			payment_url,
			payment_method,
			transaction_id,
			paid_at

		FROM orders

		WHERE user_id = ?

		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}

	for rows.Next() {

		var order Order
		var paymentToken sql.NullString
		var paymentURL sql.NullString
		var paymentMethod sql.NullString
		var transactionID sql.NullString
		var paidAt sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&paymentToken,
			&paymentURL,
			&paymentMethod,
			&transactionID,
			&paidAt,
		)
		if err != nil {
			return nil, err
		}

		if paymentToken.Valid {
			order.PaymentToken = paymentToken.String
		}
		if paymentURL.Valid {
			order.PaymentURL = paymentURL.String
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		if transactionID.Valid {
			order.TransactionID = transactionID.String
		}
		if paidAt.Valid {
			paidAtValue := paidAt.Time
			order.PaidAt = &paidAtValue
		}

		orders = append(
			orders,
			order,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) GetLibraryByUserID(
	userID int,
) ([]LibraryItem, error) {
	rows, err := r.db.Query(
		`
		SELECT
			oi.id,
			oi.order_id,
			oi.book_id,
			oi.title,
			oi.author,
			b.category,
			b.cover_url,
			b.file_path,
			o.paid_at

		FROM order_items oi
		JOIN orders o
			ON o.id = oi.order_id
		JOIN books b
			ON b.id = oi.book_id

		WHERE o.user_id = ?
			AND o.status = 'paid'

		ORDER BY o.paid_at DESC, oi.id DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	library := make([]LibraryItem, 0)

	for rows.Next() {
		var item LibraryItem
		var coverURL sql.NullString
		var filePath sql.NullString
		var paidAt sql.NullTime

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.BookID,
			&item.Title,
			&item.Author,
			&item.Category,
			&coverURL,
			&filePath,
			&paidAt,
		)
		if err != nil {
			return nil, err
		}

		if coverURL.Valid {
			item.CoverURL = coverURL.String
		}
		if filePath.Valid {
			item.FilePath = filePath.String
		}
		if paidAt.Valid {
			item.PaidAt = paidAt.Time
		}

		library = append(library, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return library, nil
}

func (r *Repository) GetOrderByID(
	userID int,
	orderID int,
) (*Order, error) {

	var order Order
	var paymentToken sql.NullString
	var paymentURL sql.NullString
	var paymentMethod sql.NullString
	var transactionID sql.NullString
	var paidAt sql.NullTime

	err := r.db.QueryRow(
		`
		SELECT
			id,
			user_id,
			total_amount,
			status,
			created_at,
			payment_token,
			payment_url,
			payment_method,
			transaction_id,
			paid_at

		FROM orders

		WHERE
			id = ?
			AND user_id = ?
		`,
		orderID,
		userID,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&paymentToken,
		&paymentURL,
		&paymentMethod,
		&transactionID,
		&paidAt,
	)

	if err != nil {
		return nil, err
	}

	if paymentToken.Valid {
		order.PaymentToken = paymentToken.String
	}
	if paymentURL.Valid {
		order.PaymentURL = paymentURL.String
	}
	if paymentMethod.Valid {
		order.PaymentMethod = paymentMethod.String
	}
	if transactionID.Valid {
		order.TransactionID = transactionID.String
	}
	if paidAt.Valid {
		paidAtValue := paidAt.Time
		order.PaidAt = &paidAtValue
	}

	rows, err := r.db.Query(
		`
		SELECT
			oi.id,
			oi.order_id,
			oi.book_id,
			oi.title,
			oi.author,
			oi.price,
			oi.quantity,
			oi.subtotal,
			b.file_path

		FROM order_items oi
		LEFT JOIN books b
			ON b.id = oi.book_id

		WHERE oi.order_id = ?
		`,
		order.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	order.Items = []OrderItem{}

	for rows.Next() {

		var item OrderItem
		var filePath sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.BookID,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.Quantity,
			&item.Subtotal,
			&filePath,
		)
		if err != nil {
			return nil, err
		}

		if order.Status == "paid" && filePath.Valid {
			item.FilePath = filePath.String
		}

		order.Items = append(
			order.Items,
			item,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}
