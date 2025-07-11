package repository

import (
	"context"
	"log/slog"

	"github.com/DENFNC/web-test/internal/domain"
	"github.com/DENFNC/web-test/internal/models"
	"github.com/DENFNC/web-test/internal/utils/dbutils"
	"github.com/DENFNC/web-test/internal/utils/mapping"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	*slog.Logger
	*goqu.DialectWrapper
	*pgxpool.Pool
}

func NewDocumentRepository(log *slog.Logger, pool *pgxpool.Pool) *DocumentRepository {
	dialect := goqu.Dialect("postgres")
	return &DocumentRepository{
		Logger:         log,
		DialectWrapper: &dialect,
		Pool:           pool,
	}
}

func (repo *DocumentRepository) SaveDocument(ctx context.Context, doc *domain.Document) (string, error) {
	var mdlDoc models.Document
	if err := mapping.MapStructModel(doc, &mdlDoc); err != nil {
		return "", err
	}

	var id string
	err := dbutils.WithTransaction(ctx, repo.Pool, func(tx pgx.Tx) error {
		stmt, args, err := repo.DialectWrapper.
			Insert("documents").
			Returning("id").
			Rows(mdlDoc).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, stmt, args...).Scan(&id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (repo *DocumentRepository) AddDocumentAccess(ctx context.Context, documentID, userID string) error {
	stmt, args, err := repo.DialectWrapper.
		Insert("document_access").
		Rows(goqu.Record{"document_id": documentID, "user_id": userID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = repo.Pool.Exec(ctx, stmt, args...)
	return err
}
