package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID        `db:"id"`
	Login     pgtype.Text        `db:"login"`
	Password  pgtype.Text        `db:"password_hash"`
	CreatedAt pgtype.Timestamptz `db:"created_at" goqu:"omitempty"`
}

type UserCredentials struct {
	Password pgtype.Text `db:"password_hash"`
	Token    pgtype.Text `db:"token"`
}
