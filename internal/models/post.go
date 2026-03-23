package models

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID string
	Title string
	Slug string
	Excerpt string
	Body string
	CoverImage string
	Tags pq.StringArray
	Published bool
	Visibility string
	ReadingTime int
	SeriesId *string
	SeriesPos *int
	CreatedAt time.Time
	UpdatedAt time.Time

	SeriesName string
	SeriesSlug string
}