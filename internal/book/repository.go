package book

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(
	search string,
	category string,
) ([]Book, error) {
	query := `
		SELECT
			id,
			title,
			author,
			category,
			description,
			price,
			cover_url,
			file_path,
			created_at
		FROM books
		WHERE 1 = 1
	`

	var args []any

	if search != "" {
		query += `
			AND (
				title LIKE ?
				OR author LIKE ?
			)
		`

		search = "%" + search + "%"

		args = append(
			args,
			search,
			search,
		)
	}

	if category != "" {
		query += `
			AND category = ?
		`

		args = append(
			args,
			category,
		)
	}

	query += `
		ORDER BY id DESC
	`

	rows, err := r.db.Query(
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := make([]Book, 0)

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Category,
			&book.Description,
			&book.Price,
			&book.CoverURL,
			&book.FilePath,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(
			books,
			book,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *Repository) GetByID(
	id int,
) (*Book, error) {
	query := `
		SELECT
			id,
			title,
			author,
			category,
			description,
			price,
			cover_url,
			file_path,
			created_at
		FROM books
		WHERE id = ?
	`

	var book Book

	err := r.db.QueryRow(
		query,
		id,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Category,
		&book.Description,
		&book.Price,
		&book.CoverURL,
		&book.FilePath,
		&book.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &book, nil
}
