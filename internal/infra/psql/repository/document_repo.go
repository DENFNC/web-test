package repository

import (
	"context"
	"fmt"
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

func (repo *DocumentRepository) GetDocumentByID(ctx context.Context, id string) (*domain.Document, error) {
	stmt, args, err := repo.DialectWrapper.
		Select("id", "file_name", "mime_type", "has_file", "is_public", "owner_id", "created_at").
		From("documents").
		Where(goqu.Ex{"id": id}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}
	row := repo.Pool.QueryRow(ctx, stmt, args...)
	var mdlDoc models.Document
	if err := row.Scan(
		&mdlDoc.ID,
		&mdlDoc.FileName,
		&mdlDoc.MimeType,
		&mdlDoc.HasFile,
		&mdlDoc.IsPublic,
		&mdlDoc.OwnerID,
		&mdlDoc.CreatedAt,
	); err != nil {
		return nil, err
	}

	var doc domain.Document
	if err := mapping.MapStructModelToDomain(&mdlDoc, &doc); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &doc, nil
}

func (repo *DocumentRepository) HasDocumentAccess(ctx context.Context, documentID, userID string) (bool, error) {
	stmt, args, err := repo.DialectWrapper.
		Select("1").
		From("document_access").
		Where(goqu.Ex{"document_id": documentID, "user_id": userID}).
		Limit(1).
		Prepared(true).
		ToSQL()
	if err != nil {
		return false, err
	}
	row := repo.Pool.QueryRow(ctx, stmt, args...)
	var dummy int
	if err := row.Scan(&dummy); err != nil {
		return false, nil
	}
	return true, nil
}

func (repo *DocumentRepository) DeleteDocument(ctx context.Context, id string) error {
	stmt, args, err := repo.DialectWrapper.
		Delete("documents").
		Where(goqu.Ex{"id": id}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = repo.Pool.Exec(ctx, stmt, args...)
	return err
}
