package cart

import (
	"database/sql"
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

func (r *Repository) GetOrCreateCart(
	userID int,
) (*Cart, error) {
	var cart Cart

	err := r.db.QueryRow(
		`
		SELECT
			id,
			user_id
		FROM carts
		WHERE user_id = ?
		`,
		userID,
	).Scan(
		&cart.ID,
		&cart.UserID,
	)

	if err == sql.ErrNoRows {

		result, err := r.db.Exec(
			`
			INSERT INTO carts (
				user_id
			)
			VALUES (?)
			`,
			userID,
		)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		cart.ID = int(id)
		cart.UserID = userID

		return &cart, nil
	}

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *Repository) AddItem(
	cartID int,
	bookID int,
) error {
	_, err := r.db.Exec(
		`
		INSERT INTO cart_items (
			cart_id,
			book_id,
			quantity
		)
		VALUES (?, ?, 1)

		ON CONFLICT(cart_id, book_id)
		DO UPDATE SET
			quantity = quantity + 1
		`,
		cartID,
		bookID,
	)

	return err
}

func (r *Repository) GetCart(
	userID int,
) (*Cart, error) {
	cart, err := r.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		`
		SELECT
			cart_items.id,
			books.id,
			books.title,
			books.author,
			books.price,
			books.cover_url,
			cart_items.quantity

		FROM cart_items

		JOIN books
			ON books.id = cart_items.book_id

		WHERE cart_items.cart_id = ?
		`,
		cart.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []CartItem{}

	for rows.Next() {
		var item CartItem

		err := rows.Scan(
			&item.ID,
			&item.BookID,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.CoverURL,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		items = append(
			items,
			item,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	cart.Items = items

	return cart, nil
}

func (r *Repository) UpdateItemQuantity(
	cartID int,
	itemID int,
	quantity int,
) error {
	_, err := r.db.Exec(
		`
		UPDATE cart_items
		SET quantity = ?
		WHERE
			id = ?
			AND cart_id = ?
		`,
		quantity,
		itemID,
		cartID,
	)

	return err
}

func (r *Repository) DeleteItem(
	cartID int,
	itemID int,
) error {
	_, err := r.db.Exec(
		`
		DELETE FROM cart_items
		WHERE
			id = ?
			AND cart_id = ?
		`,
		itemID,
		cartID,
	)

	return err
}
