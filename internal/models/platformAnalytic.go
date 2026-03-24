package models

import "time"

type PlatformAnalytics struct {
	ID            string    `json:"id"`
	PublicationID string 	`json:"publicationId"`
	Views         int 		`json:"views"`
	Reactions     int		`json:"reactions"`
	Comments      int		`json:"comments"`
	Bookmarks     int		`json:"bookmarks"`
	FetchedAt     time.Time	`json:"fetchedAt"`
}