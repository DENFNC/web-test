package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"

	"github.com/DENFNC/web-test/internal/model"
	"github.com/DENFNC/web-test/internal/repository"
)

type UserService struct {
	repo       *repository.MemoryRepo
	adminToken string
}

func NewUserService(repo *repository.MemoryRepo, adminToken string) *UserService {
	return &UserService{repo: repo, adminToken: adminToken}
}

// Регистрация пользователя
func (s *UserService) Register(adminToken, login, pswd string) error {
	if adminToken != s.adminToken {
		return errors.New("forbidden")
	}
	if !isValidLogin(login) {
		return errors.New("invalid login")
	}
	if !isValidPassword(pswd) {
		return errors.New("invalid password")
	}
	if _, exists := s.repo.GetUser(login); exists {
		return errors.New("user exists")
	}
	hash := hashPassword(pswd)
	user := &model.User{Login: login, PasswordHash: hash}
	s.repo.AddUser(user)
	return nil
}

// Аутентификация
func (s *UserService) Auth(login, pswd string) (string, error) {
	user, ok := s.repo.GetUser(login)
	if !ok || user.PasswordHash != hashPassword(pswd) {
		return "", errors.New("unauthorized")
	}
	token := genToken()
	s.repo.AddSession(token, login)
	return token, nil
}

// Проверка токена
func (s *UserService) CheckToken(token string) (string, bool) {
	sess, ok := s.repo.GetSession(token)
	if !ok {
		return "", false
	}
	return sess.Login, true
}

// Завершение сессии
func (s *UserService) Logout(token string) {
	s.repo.DeleteSession(token)
}

// --- Вспомогательные функции ---
func isValidLogin(login string) bool {
	if len(login) < 8 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, login)
	return matched
}

func isValidPassword(pswd string) bool {
	if len(pswd) < 8 {
		return false
	}
	var upper, lower, digit, special bool
	for _, c := range pswd {
		switch {
		case c >= 'A' && c <= 'Z':
			upper = true
		case c >= 'a' && c <= 'z':
			lower = true
		case c >= '0' && c <= '9':
			digit = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>/?`~", c):
			special = true
		}
	}
	return upper && lower && digit && special
}

func hashPassword(pswd string) string {
	h := sha256.Sum256([]byte(pswd))
	return base64.StdEncoding.EncodeToString(h[:])
}

func genToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
