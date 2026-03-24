package models

import "time"

type Series struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"createdAt"`

	PostCount int
}