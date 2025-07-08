package repository

import (
	"time"

	"github.com/DENFNC/web-test/internal/model"
)

type Repo interface {
	// Пользователи
	GetUser(login string) (*model.User, bool)
	AddUser(user *model.User)
	// Сессии
	AddSession(token, login string)
	GetSession(token string) (*model.Session, bool)
	DeleteSession(token string)
	// Документы
	ListDocuments() []*model.Document
	GetDocument(id string) (*model.Document, bool)
	AddDocument(doc *model.Document)
	DeleteDocument(id string)
}

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
	Invalidate(keys ...string)
}
