package store

import (
	"context"
	"database/sql"
	"fmt"

	"origin.me/internal/models"
)

func (s *Store) GetSeriesBySlug(ctx context.Context, slug string) (*models.Series, error) {
	sr := &models.Series{}
	err := s.db.QueryRowContext(ctx, `
		SELECT 
			s."id",
			s."name",
			s."slug",
			COALESCE(s."description", ''),
			s."createdAt",
			COUNT(p."id") AS "postCount"
		FROM series s
		LEFT JOIN posts p ON s."id" = p."seriesId" and p."published" = true
		WHERE s."slug" = $1
		GROUP BY s."id"
	`, slug).Scan(
		&sr.ID,
		&sr.Name,
		&sr.Slug,
		&sr.Description,
		&sr.CreatedAt,
		&sr.PostCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return sr, err
}

func (s *Store) GetSeriesByID(ctx context.Context, id string) (*models.Series, error) {
	sr := &models.Series{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, slug, COALESCE(description,''), "createdAt"
		FROM series WHERE id = $1`, id,
	).Scan(&sr.ID, &sr.Name, &sr.Slug, &sr.Description, &sr.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return sr, err
}

func (s *Store) ListSeries(ctx context.Context) ([]*models.Series, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			s.id, s.name, s.slug,
			COALESCE(s.description, ''),
			s."createdAt",
			COUNT(p.id) as post_count
		FROM series s
		LEFT JOIN posts p ON p."seriesId" = s.id AND p.published = true
		GROUP BY s.id
		ORDER BY s."createdAt" DESC`)
	if err != nil {
		return nil, fmt.Errorf("list series: %w", err)
	}
	defer rows.Close()

	var out []*models.Series
	for rows.Next() {
		sr := &models.Series{}
		if err := rows.Scan(
			&sr.ID, &sr.Name, &sr.Slug,
			&sr.Description, &sr.CreatedAt, &sr.PostCount,
		); err != nil {
			return nil, err
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}

func (s *Store) CreateSeries(ctx context.Context, sr *models.Series) error {
	return s.db.QueryRowContext(ctx, `
		INSERT INTO series (name, slug, description)
		VALUES ($1, $2, $3)
		RETURNING id, "createdAt"`,
		sr.Name, sr.Slug, sr.Description,
	).Scan(&sr.ID, &sr.CreatedAt)
}

func (s *Store) UpdateSeries(ctx context.Context, sr *models.Series) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE series SET name=$1, slug=$2, description=$3 WHERE id=$4`,
		sr.Name, sr.Slug, sr.Description, sr.ID)
	return err
}

func (s *Store) DeleteSeries(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM series WHERE id=$1`, id)
	return err
}