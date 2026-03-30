package store

import (
	"context"
	"fmt"
	"log"

	"origin.me/internal/models"
)

func (s *Store) GetLatestAnalytics(ctx context.Context, postID string) ([]*models.PostPublication, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			pp.id,
			pp."postId",
			pp.platform,
			COALESCE(pp."platformPostId", ''),
			COALESCE(pp."platformUrl", ''),
			pp."canonicalUrl",
			pp."syncStatus",
			pp."syncedAt",
			pp."createdAt",
			COALESCE(pa.views, 0),
			COALESCE(pa.reactions, 0),
			COALESCE(pa.comments, 0),
			COALESCE(pa.bookmarks, 0),
			COALESCE(pa."fetchedAt", NOW())
		FROM post_publications pp
		LEFT JOIN LATERAL (
			SELECT views, reactions, comments, bookmarks, "fetchedAt"
			FROM platform_analytics
			WHERE "publicationId" = pp.id
			ORDER BY "fetchedAt" DESC
			LIMIT 1
		) pa ON true
		WHERE pp."postId" = $1
		ORDER BY pp."createdAt" ASC`,
		postID,
	)
	if err != nil {
		return nil, fmt.Errorf("get latest analytics: %w", err)
	}
	defer rows.Close()

	var out []*models.PostPublication
	for rows.Next() {
		p := &models.PostPublication{}
		a := &models.PlatformAnalytics{}

		if err := rows.Scan(
			&p.ID,
			&p.PostId,
			&p.Platform,
			&p.PlatformPostID,
			&p.PlatformURL,
			&p.CanonicalURL,
			&p.SyncStatus,
			&p.SyncedAt,
			&p.CreatedAt,
			&a.Views,
			&a.Reactions,
			&a.Comments,
			&a.Bookmarks,
			&a.FetchedAt,
		); err != nil {
			return nil, fmt.Errorf("scan publication: %w", err)
		}

		p.Analytics = a
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetTotalAnalytics(ctx context.Context, postID string) (views, reactions, comments int, err error) {
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(pa.views), 0),
			COALESCE(SUM(pa.reactions), 0),
			COALESCE(SUM(pa.comments), 0)
		FROM post_publications pp
		LEFT JOIN LATERAL (
			SELECT views, reactions, comments
			FROM platform_analytics
			WHERE "publicationId" = pp.id
			ORDER BY "fetchedAt" DESC
			LIMIT 1
		) pa ON true
		WHERE pp."postId" = $1`,
		postID,
	).Scan(&views, &reactions, &comments)

	if err != nil {
		err = fmt.Errorf("get total analytics: %w", err)
	}
	return
}

func (s *Store) GetPostSiteViews(ctx context.Context, postID string) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(views), 0)
		FROM site_analytics
		WHERE "postId" = $1`,
		postID,
	).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("get post site views: %w", err)
	}
	return total, nil
}

func (s *Store) GetSiteViewsLast14Days(ctx context.Context) ([]*models.SiteAnalytics, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			'all'           AS "postId",
			SUM(views)      AS views,
			SUM("uniqueViews") AS "uniqueViews",
			0               AS "avgReadTimes",
			0               AS "avgScrollDepth",
			date
		FROM site_analytics
		WHERE date >= CURRENT_DATE - INTERVAL '14 days'
		GROUP BY date
		ORDER BY date ASC`)
	if err != nil {
		return nil, fmt.Errorf("get site views 14 days: %w", err)
	}
	defer rows.Close()

	var out []*models.SiteAnalytics
	for rows.Next() {
		a := &models.SiteAnalytics{}
		if err := rows.Scan(
			&a.PostID,
			&a.Views,
			&a.UniqueViews,
			&a.AvgReadingTime,
			&a.AvgScrollDepth,
			&a.Date,
		); err != nil {
			return nil, fmt.Errorf("scan site analytics: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) GetSiteAnalytics(ctx context.Context, postID string, days int) ([]*models.SiteAnalytics, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			"postId", views, "uniqueViews",
			"avgReadTimes", "avgScrollDepth", date
		FROM site_analytics
		WHERE "postId" = $1
		  AND date >= CURRENT_DATE - ($2 * INTERVAL '1 day')
		ORDER BY date ASC`,
		postID, days,
	)
	if err != nil {
		return nil, fmt.Errorf("get site analytics: %w", err)
	}
	defer rows.Close()

	var out []*models.SiteAnalytics
	for rows.Next() {
		a := &models.SiteAnalytics{}
		if err := rows.Scan(
			&a.PostID,
			&a.Views,
			&a.UniqueViews,
			&a.AvgReadingTime,
			&a.AvgScrollDepth,
			&a.Date,
		); err != nil {
			return nil, fmt.Errorf("scan site analytics: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (h *Store) TrackSiteView(ctx context.Context, postID string) {
	_, err := h.db.ExecContext(ctx, `
		INSERT INTO site_analytics ("postId", views, "uniqueViews")
		VALUES ($1, 1, 1)
		ON CONFLICT ("postId", date) DO UPDATE SET
			views = site_analytics.views + 1,
			"uniqueViews" = site_analytics."uniqueViews" + 1`,
		postID,
	)
	if err != nil {
		log.Printf("track site view: %v", err)
	}
}

