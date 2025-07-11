package domain

import "time"

type Document struct {
	ID        string
	FileName  string
	MimeType  string
	HasFile   bool
	IsPublic  bool
	OwnerID   string
	CreatedAt time.Time
}
