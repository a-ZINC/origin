package models

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID string  `db:"id"`
	Title string `db:"title"`
	Slug string `db:"slug"`
	Excerpt string `db:"excerpt"`
	Body string `db:"body"`
	CoverImage string `db:"coverImage"`
	Tags pq.StringArray `db:"tags"`
	Published bool `db:"published"`
	Visibility string `db:"visibility"`
	ReadingTime int `db:"readingTime"`
	SeriesId *string `db:"seriesId"`
	SeriesPos *int `db:"seriesPos"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`

	SeriesName string
	SeriesSlug string

	HTMLBody string
	PrevPost *PostNav
    NextPost *PostNav
}

type PostNav struct {
    Title string
    Slug  string
}