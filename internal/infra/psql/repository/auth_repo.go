package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/DENFNC/web-test/internal/domain"
	"github.com/DENFNC/web-test/internal/models"
	"github.com/DENFNC/web-test/internal/utils/dbutils"
	"github.com/DENFNC/web-test/internal/utils/mapping"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	*slog.Logger
	*goqu.DialectWrapper
	*pgxpool.Pool
}

func NewAuthRepository(
	log *slog.Logger,
	pool *pgxpool.Pool,
) *AuthRepository {
	dialect := goqu.Dialect("postgres")

	return &AuthRepository{
		Logger:         log,
		DialectWrapper: &dialect,
		Pool:           pool,
	}
}

func (repo *AuthRepository) GetByID(ctx context.Context, id string) (*domain.UserCredentials, error) {
	stmt, args, err := repo.DialectWrapper.
		Select(
			goqu.I("users.password_hash"),
			goqu.I("auth_tokens.token"),
		).
		From("users").
		LeftJoin(
			goqu.T("auth_tokens"),
			goqu.On(goqu.Ex{
				"users.id": goqu.I("auth_tokens.user_id"),
			}),
		).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := repo.Pool.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	creds, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserCredentials])
	if err != nil {
		return nil, err
	}

	var userCreds domain.UserCredentials
	fmt.Println(userCreds)
	if err := mapping.MapStructModelToDomain(creds, &userCreds); err != nil {
		return nil, err
	}

	return &userCreds, nil
}

func (repo *AuthRepository) SaveUser(ctx context.Context, user *domain.User) (string, error) {
	var mdlUser models.User
	if err := mapping.MapStructModel(user, &mdlUser); err != nil {
		return "", err
	}

	var login string
	err := dbutils.WithTransaction(ctx, repo.Pool, func(tx pgx.Tx) error {
		var err error

		login, err = repo.insertUser(ctx, tx, &mdlUser)
		if err != nil {
			return err
		}

		token, err := generateToken(128)
		if err != nil {
			return err
		}

		authToken := models.AuthToken{
			UserID: mdlUser.ID,
			Token: pgtype.Text{
				String: token,
				Valid:  true,
			},
		}

		if err := repo.insertAuthToken(ctx, tx, &authToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return login, nil
}

func (repo *AuthRepository) insertUser(ctx context.Context, tx pgx.Tx, user *models.User) (string, error) {
	stmt, args, err := repo.DialectWrapper.
		Insert("users").
		Returning("login").
		Rows(user).
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}

	var login string
	if err := tx.QueryRow(ctx, stmt, args...).Scan(&login); err != nil {
		return "", err
	}

	return login, nil
}

func (repo *AuthRepository) insertAuthToken(ctx context.Context, tx pgx.Tx, token *models.AuthToken) error {
	stmt, args, err := repo.DialectWrapper.
		Insert("auth_tokens").
		Rows(token).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, stmt, args...)
	return err
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (repo *AuthRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	stmt, args, err := repo.DialectWrapper.
		Select("user_id").
		From("auth_tokens").
		Where(goqu.Ex{"token": token, "is_revoked": false}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}
	row := repo.Pool.QueryRow(ctx, stmt, args...)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (repo *AuthRepository) GetUserIDByLogin(ctx context.Context, login string) (string, error) {
	stmt, args, err := repo.DialectWrapper.
		Select("id").
		From("users").
		Where(goqu.Ex{"login": login}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}
	row := repo.Pool.QueryRow(ctx, stmt, args...)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}
