package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Document struct {
	ID        pgtype.UUID        `db:"id"`
	FileName  pgtype.Text        `db:"file_name"`
	MimeType  pgtype.Text        `db:"mime_type"`
	HasFile   pgtype.Bool        `db:"has_file"`
	IsPublic  pgtype.Bool        `db:"is_public"`
	OwnerID   pgtype.UUID        `db:"owner_id"`
	CreatedAt pgtype.Timestamptz `db:"created_at" goqu:"omitempty"`
}
