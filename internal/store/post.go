package store

import (
	"context"
	"database/sql"
	"fmt"

	"origin.me/internal/models"
)

func scanPost(row interface {
	Scan(...interface{}) error
}, p *models.Post) error {
	return row.Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Excerpt,
		&p.Body,
		&p.CoverImage,
		&p.Tags,
		&p.Published,
		&p.Visibility,
		&p.ReadingTime,
		&p.SeriesId,
		&p.SeriesPos,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.SeriesName,
		&p.SeriesSlug,
	)
}

func (s *Store) GetPostBySlug(ctx context.Context, slug string) (*models.Post, error) {
	p := &models.Post{}

	err := s.db.QueryRowContext(ctx, `
		SELECT
			p."id",
			p."title",
			p."slug",
			COALESCE(p."excerpt", ''),
			p."body",
			COALESCE(p."coverImage", ''),
			p."tags",
			p."published",
			p."visibility",
			p."readingTime",
			p."seriesId",
			p."seriesPos",
			p."createdAt",
			p."updatedAt",
			COALESCE(s."name", ''),
			COALESCE(s."slug", '')
		FROM posts p
		LEFT JOIN series s ON p."seriesId" = s."id"
		WHERE p.slug = $1
	`, slug).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Excerpt,
		&p.Body,
		&p.CoverImage,
		&p.Tags,
		&p.Published,
		&p.Visibility,
		&p.ReadingTime,
		&p.SeriesId,
		&p.SeriesPos,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.SeriesName,
		&p.SeriesSlug,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get post by slug: %w", err)
	}

	return p, nil
}

func (s *Store) ListPosts(ctx context.Context, includePrivate bool) ([]*models.Post, error) {
	query := `
		SELECT
			p."id",
			p."title",
			p."slug",
			COALESCE(p."excerpt", ''),
			p."body",
			COALESCE(p."coverImage", ''),
			p."tags",
			p."published",
			p."visibility",
			p."readingTime",
			p."seriesId",
			p."seriesPos",
			p."createdAt",
			p."updatedAt",
			COALESCE(s."name", ''),
			COALESCE(s."slug", '')
		FROM posts p
		LEFT JOIN series s ON p."seriesId" = s."id"
		WHERE p."published" = true
	`

	if !includePrivate {
		query += " AND p.visibility = 'public'"
	}

	query += ` ORDER BY p."createdAt" DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Excerpt,
			&p.Body,
			&p.CoverImage,
			&p.Tags,
			&p.Published,
			&p.Visibility,
			&p.ReadingTime,
			&p.SeriesId,
			&p.SeriesPos,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.SeriesName,
			&p.SeriesSlug,
		)
		if err != nil {
			return nil, fmt.Errorf("list posts: %w", err)
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}

	return posts, nil
}

func (s *Store) ListPostsByTag(ctx context.Context, tag string, includePrivate bool) ([]*models.Post, error) {
	query := `
		SELECT
			p."id",
			p."title",
			p."slug",
			COALESCE(p."excerpt", ''),
			p."body",
			COALESCE(p."coverImage", ''),
			p."tags",
			p."published",
			p."visibility",
			p."readingTime",
			p."seriesId",
			p."seriesPos",
			p."createdAt",
			p."updatedAt",
			COALESCE(s."name", ''),
			COALESCE(s."slug", '')
		FROM posts p
		LEFT JOIN series s ON p."seriesId" = s."id"
		WHERE p."published" = true
		AND $1 = ANY(p."tags")
	`

	if !includePrivate {
		query += " AND p.visibility = 'public'"
	}

	query += ` ORDER BY p."createdAt" DESC`

	rows, err := s.db.QueryContext(ctx, query, tag)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Excerpt,
			&p.Body,
			&p.CoverImage,
			&p.Tags,
			&p.Published,
			&p.Visibility,
			&p.ReadingTime,
			&p.SeriesId,
			&p.SeriesPos,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.SeriesName,
			&p.SeriesSlug,
		)
		if err != nil {
			return nil, fmt.Errorf("list posts: %w", err)
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}

	return posts, nil
}