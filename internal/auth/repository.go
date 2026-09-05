package auth

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

func (r *Repository) CreateUser(
	name string,
	email string,
	password string,
) (*User, error) {
	result, err := r.db.Exec(
		`
		INSERT INTO users (
			name,
			email,
			password
		)
		VALUES (?, ?, ?)
		`,
		name,
		email,
		password,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       int(id),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "user",
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(
	email string,
) (*User, error) {
	var user User

	err := r.db.QueryRow(
		`
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at
		FROM users
		WHERE email = ?
		`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateSession(
	userID int,
	token string,
	expiresAt string,
) error {
	_, err := r.db.Exec(
		`
		INSERT INTO sessions (
			user_id,
			token,
			expires_at
		)
		VALUES (?, ?, ?)
		`,
		userID,
		token,
		expiresAt,
	)

	return err
}

func (r *Repository) DeleteSession(
	token string,
) error {
	_, err := r.db.Exec(
		`
		DELETE FROM sessions
		WHERE token = ?
		`,
		token,
	)

	return err
}

func (r *Repository) GetUserBySession(
	token string,
) (*User, error) {
	var user User

	err := r.db.QueryRow(
		`
		SELECT
			users.id,
			users.name,
			users.email,
			users.role,
			users.created_at
		FROM sessions

		JOIN users
			ON users.id = sessions.user_id

		WHERE
			sessions.token = ?
			AND sessions.expires_at > CURRENT_TIMESTAMP
		`,
		token,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
