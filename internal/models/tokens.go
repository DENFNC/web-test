package models

import "github.com/jackc/pgx/v5/pgtype"

type AuthToken struct {
	UserID    pgtype.UUID        `db:"user_id"`
	Token     pgtype.Text        `db:"token"`
	IsRevoked pgtype.Bool        `db:"is_revoked" goqu:"omitempty"`
	CreatedAt pgtype.Timestamptz `db:"created_at" goqu:"omitempty"`
}
