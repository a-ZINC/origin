package store

import (
	"context"
	"database/sql"

	"origin.me/internal/models"
)

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, email, username, password, "createdAt" FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (s *Store) CreateUser(ctx context.Context, email, username, hashedPassword string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (email, username, password) VALUES ($1, $2, $3)`,
		email, username, hashedPassword,
	)
	return err
}