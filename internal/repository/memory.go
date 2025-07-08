package repository

import (
	"sync"
	"time"

	"github.com/DENFNC/web-test/internal/model"
)

type MemoryRepo struct {
	mu        sync.RWMutex
	Documents map[string]*model.Document
	Users     map[string]*model.User
	Sessions  map[string]*model.Session // token -> session
}

func NewMemoryRepo() *MemoryRepo {
	repo := &MemoryRepo{
		Documents: make(map[string]*model.Document),
		Users:     make(map[string]*model.User),
		Sessions:  make(map[string]*model.Session),
	}
	// Тестовые документы
	repo.Documents["doc1"] = &model.Document{
		ID:      "doc1",
		Name:    "photo.jpg",
		Mime:    "image/jpg",
		File:    true,
		Public:  false,
		Created: time.Now(),
		Grant:   []string{"login1", "login2"},
		Owner:   "testuser",
	}
	repo.Documents["doc2"] = &model.Document{
		ID:      "doc2",
		Name:    "text.txt",
		Mime:    "text/plain",
		File:    false,
		Public:  true,
		Created: time.Now(),
		Grant:   []string{},
		Owner:   "testuser",
		JSON:    map[string]any{"content": "hello"},
	}
	return repo
}

// Документы
func (r *MemoryRepo) ListDocuments() []*model.Document {
	r.mu.RLock()
	defer r.mu.RUnlock()
	docs := make([]*model.Document, 0, len(r.Documents))
	for _, d := range r.Documents {
		docs = append(docs, d)
	}
	return docs
}

func (r *MemoryRepo) GetDocument(id string) (*model.Document, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.Documents[id]
	return d, ok
}

func (r *MemoryRepo) AddDocument(doc *model.Document) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Documents[doc.ID] = doc
}

func (r *MemoryRepo) DeleteDocument(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Documents, id)
}

// Пользователи (для будущих этапов)
func (r *MemoryRepo) GetUser(login string) (*model.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.Users[login]
	return u, ok
}

func (r *MemoryRepo) AddUser(user *model.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Users[user.Login] = user
}

// Сессии
func (r *MemoryRepo) AddSession(token, login string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Sessions[token] = &model.Session{Token: token, Login: login}
}

func (r *MemoryRepo) GetSession(token string) (*model.Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.Sessions[token]
	return s, ok
}

func (r *MemoryRepo) DeleteSession(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Sessions, token)
}
