package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

const postSelect = `
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
`

func (s *Store) GetPostBySlug(ctx context.Context, slug string) (*models.Post, error) {
	p := &models.Post{}

	row := s.db.QueryRowContext(ctx, postSelect+`
		WHERE p.slug = $1
	`, slug)

	err := scanPost(row, p)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get post by slug: %w", err)
	}

	return p, nil
}

func (s *Store) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	p := &models.Post{}
	row := s.db.QueryRowContext(ctx,
		postSelect+` WHERE p.id = $1`, id)

	if err := scanPost(row, p); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	return p, nil
}

func (s *Store) ListPosts(ctx context.Context, includePrivate bool) ([]*models.Post, error) {
	query := postSelect + `WHERE p."published" = true`

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

		err := scanPost(rows, p)
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
	query := postSelect + `
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
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}

		err := scanPost(rows, p)
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

func (s *Store) ListAllPosts(ctx context.Context) ([]*models.Post, error) {
	query := postSelect + `ORDER BY p."createdAt" DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}

		if err := scanPost(rows, p); err != nil {
			return nil, fmt.Errorf("list posts: %w", err)
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (s *Store) CreatePost(ctx context.Context, p *models.Post) error {
	p.ReadingTime = models.CalcReadingTime(p.Body)

	var publishedAt interface{}
	if p.Published {
		publishedAt = time.Now()
	}
	query := `
		INSERT INTO posts (
			"title",
			"slug",
			"excerpt",
			"body",
			"coverImage",
			"tags",
			"published",
			"visibility",
			"readingTime",
			"seriesId",
			"seriesPos",
			"createdAt",
			"updatedAt"
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
		)
		returning id, "createdAt", "updatedAt"
	`

	return s.db.QueryRowContext(ctx, query,
		p.Title,
		p.Slug,
		p.Excerpt,
		p.Body,
		p.CoverImage,
		p.Tags,
		p.Published,
		p.Visibility,
		p.ReadingTime,
		p.SeriesId,
		p.SeriesPos,
		publishedAt,
		time.Now(),
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (s *Store) UpdatePost(ctx context.Context, p *models.Post) error {
	p.ReadingTime = models.CalcReadingTime(p.Body)

	query := `
		UPDATE posts SET
			"title" = $1,
			"slug" = $2,
			"excerpt" = $3,
			"body" = $4,
			"coverImage" = $5,
			"tags" = $6,
			"published" = $7,
			"visibility" = $8,
			"readingTime" = $9,
			"seriesId" = $10,
			"seriesPos" = $11,
			"updatedAt" = NOW()
		WHERE id = $12
	`

	_, err := s.db.ExecContext(ctx, query,
		p.Title,
		p.Slug,
		p.Excerpt,
		p.Body,
		p.CoverImage,
		p.Tags,
		p.Published,
		p.Visibility,
		p.ReadingTime,
		p.SeriesId,
		p.SeriesPos,
		p.ID,
	)

	return err
}

func (s *Store) DeletePost(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}