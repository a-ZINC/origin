package models

import "time"

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}