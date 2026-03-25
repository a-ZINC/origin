package store

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (s *Store) SeedAdmin(ctx context.Context, email, username, password string) error {
	existing, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if existing != nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.CreateUser(ctx, email, username, string(hash))
}