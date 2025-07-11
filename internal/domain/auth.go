package domain

import "time"

type AuthToken struct {
	UserID    string
	Token     string
	IsRevoked bool
	CreatedAt time.Time
}
