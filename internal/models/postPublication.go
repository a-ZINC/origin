package models

import "time"

type PostPublication struct {
	ID string `json:"id"`
	PostId string `json:"postId"`
	Platform string `json:"platform"`
	PlatformPostID string `json:"platformPostId"`
	PlatformURL string `json:"platformURL"`
	CanonicalURL  string `json:"canonicalURL"`
	SyncStatus string `json:"syncStatus"`
	SyncedAt *time.Time `json:"syncedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Analytics *PlatformAnalytics
}