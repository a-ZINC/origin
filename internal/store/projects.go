package store

import (
	"context"
	"fmt"

	"origin.me/internal/models"
)

func (s *Store) ListProjects(ctx context.Context) ([]*models.Project, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id, name, description,
			COALESCE(url, ''), COALESCE("repoUrl", ''),
			COALESCE(language, ''),
			stars, featured, "sortOrder", "createdAt"
		FROM projects
		ORDER BY "sortOrder" ASC, stars DESC`)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var out []*models.Project
	for rows.Next() {
		pr := &models.Project{}
		if err := rows.Scan(
			&pr.ID, &pr.Name, &pr.Description,
			&pr.URL, &pr.RepoURL, &pr.Language,
			&pr.Stars, &pr.Featured, &pr.SortOrder, &pr.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, pr)
	}
	return out, rows.Err()
}

func (s *Store) CreateProject(ctx context.Context, pr *models.Project) error {
	return s.db.QueryRowContext(ctx, `
		INSERT INTO projects (name, description, url, "repoUrl", language, stars, featured, "sortOrder")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, "createdAt"`,
		pr.Name, pr.Description, pr.URL, pr.RepoURL,
		pr.Language, pr.Stars, pr.Featured, pr.SortOrder,
	).Scan(&pr.ID, &pr.CreatedAt)
}

func (s *Store) UpdateProject(ctx context.Context, pr *models.Project) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE projects SET
			name=$1, description=$2, url=$3, "repoUrl"=$4,
			language=$5, stars=$6, featured=$7, "sortOrder"=$8
		WHERE id=$9`,
		pr.Name, pr.Description, pr.URL, pr.RepoURL,
		pr.Language, pr.Stars, pr.Featured, pr.SortOrder, pr.ID)
	return err
}

func (s *Store) DeleteProject(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM projects WHERE id=$1`, id)
	return err
}