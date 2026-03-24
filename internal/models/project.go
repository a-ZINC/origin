package models

import "time"

type Project struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Language string `json:"language"`
	RepoURL string `json:"repoUrl"`
	URL string `json:"url"`
	Stars int `json:"stars"`
	Featured bool `json:"featured"`
	SortOrder int `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
}