package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"

	"github.com/DENFNC/web-test/internal/domain"
)

type UserUsecase struct {
	userRepo    domain.UserRepo
	sessionRepo domain.SessionRepo
	adminToken  string
}

func NewUserUsecase(userRepo domain.UserRepo, sessionRepo domain.SessionRepo, adminToken string) *UserUsecase {
	return &UserUsecase{userRepo, sessionRepo, adminToken}
}

func (u *UserUsecase) Register(adminToken, login, pswd string) error {
	if adminToken != u.adminToken {
		return errors.New("forbidden")
	}
	if !isValidLogin(login) {
		return errors.New("invalid login")
	}
	if !isValidPassword(pswd) {
		return errors.New("invalid password")
	}
	_, err := u.userRepo.Get(login)
	if err == nil {
		return errors.New("user exists")
	}
	hash := hashPassword(pswd)
	return u.userRepo.Add(&domain.User{Login: login, PasswordHash: hash})
}

func (u *UserUsecase) Auth(login, pswd string) (string, error) {
	user, err := u.userRepo.Get(login)
	if err != nil || user.PasswordHash != hashPassword(pswd) {
		return "", errors.New("unauthorized")
	}
	token := genToken()
	u.sessionRepo.Add(token, login)
	return token, nil
}

func (u *UserUsecase) CheckToken(token string) (string, bool) {
	sess, err := u.sessionRepo.Get(token)
	if err != nil {
		return "", false
	}
	return sess.Login, true
}

func (u *UserUsecase) Logout(token string) {
	u.sessionRepo.Delete(token)
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
