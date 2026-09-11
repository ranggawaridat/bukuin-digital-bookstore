package order

import (
	"database/sql"
	"errors"
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

	return order, nil
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
			created_at

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

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
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

func (r *Repository) GetOrderByID(
	userID int,
	orderID int,
) (*Order, error) {

	var order Order

	err := r.db.QueryRow(
		`
		SELECT
			id,
			user_id,
			total_amount,
			status,
			created_at

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
	)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		`
		SELECT
			id,
			order_id,
			book_id,
			title,
			author,
			price,
			quantity,
			subtotal

		FROM order_items

		WHERE order_id = ?
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

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.BookID,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.Quantity,
			&item.Subtotal,
		)
		if err != nil {
			return nil, err
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
