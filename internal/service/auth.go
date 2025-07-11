package service

import (
	"context"
	"log/slog"

	"github.com/DENFNC/web-test/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetByID(ctx context.Context, id string) (*domain.UserCredentials, error)
	SaveUser(ctx context.Context, user *domain.User) (string, error)
}

type AuthService struct {
	*slog.Logger
	repo AuthRepository
}

func NewAuthService(log *slog.Logger, repo AuthRepository) *AuthService {
	return &AuthService{
		Logger: log,
		repo:   repo,
	}
}

func (srv *AuthService) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	const op = "service.AuthService.CreateUser"

	log := srv.Logger.With("op", op)

	hash, err := HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	userID := uuid.New()
	user.ID = userID.String()
	user.Password = string(hash)

	login, err := srv.repo.SaveUser(ctx, user)
	if err != nil {
		log.Error(
			"Entity creation error",
			slog.String("err", err.Error()),
		)
		return "", err
	}

	return login, nil
}

func (srv *AuthService) LoginUser(ctx context.Context, user *domain.User) (string, error) {
	const op = "service.AuthService.LoginUser"

	_ = srv.Logger.With("op", op)

	data, err := srv.repo.GetByID(ctx, user.ID)
	if err != nil {
		return "", err
	}

	if !CheckPasswordHash(user.Password, data.Password) {
		return "", nil
	}

	return data.Token, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
