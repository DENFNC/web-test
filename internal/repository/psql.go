package repository

import (
	"context"

	"github.com/DENFNC/web-test/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PsqlRepo struct {
	db *pgxpool.Pool
}

func NewPsqlRepo(db *pgxpool.Pool) *PsqlRepo {
	return &PsqlRepo{db: db}
}

// --- Пользователи ---
func (r *PsqlRepo) GetUser(login string) (*model.User, bool) {
	row := r.db.QueryRow(context.Background(), "SELECT login, password_hash FROM users WHERE login=$1", login)
	var u model.User
	if err := row.Scan(&u.Login, &u.PasswordHash); err != nil {
		return nil, false
	}
	return &u, true
}

func (r *PsqlRepo) AddUser(user *model.User) {
	r.db.Exec(context.Background(), "INSERT INTO users (login, password_hash) VALUES ($1, $2)", user.Login, user.PasswordHash)
}

// --- Сессии ---
func (r *PsqlRepo) AddSession(token, login string) {
	r.db.Exec(context.Background(), "INSERT INTO sessions (token, login) VALUES ($1, $2)", token, login)
}

func (r *PsqlRepo) GetSession(token string) (*model.Session, bool) {
	row := r.db.QueryRow(context.Background(), "SELECT token, login FROM sessions WHERE token=$1", token)
	var s model.Session
	if err := row.Scan(&s.Token, &s.Login); err != nil {
		return nil, false
	}
	return &s, true
}

func (r *PsqlRepo) DeleteSession(token string) {
	r.db.Exec(context.Background(), "DELETE FROM sessions WHERE token=$1", token)
}

// --- Документы ---
func (r *PsqlRepo) ListDocuments() []*model.Document {
	rows, _ := r.db.Query(context.Background(), "SELECT id, name, mime, file, public, created, grant, owner, json FROM documents")
	var docs []*model.Document
	for rows.Next() {
		var d model.Document
		var grant []string
		var jsonData []byte
		rows.Scan(&d.ID, &d.Name, &d.Mime, &d.File, &d.Public, &d.Created, &grant, &d.Owner, &jsonData)
		d.Grant = grant
		if len(jsonData) > 0 {
			d.JSON = string(jsonData)
		}
		docs = append(docs, &d)
	}
	return docs
}

func (r *PsqlRepo) GetDocument(id string) (*model.Document, bool) {
	row := r.db.QueryRow(context.Background(), "SELECT id, name, mime, file, public, created, grant, owner, json FROM documents WHERE id=$1", id)
	var d model.Document
	var grant []string
	var jsonData []byte
	if err := row.Scan(&d.ID, &d.Name, &d.Mime, &d.File, &d.Public, &d.Created, &grant, &d.Owner, &jsonData); err != nil {
		return nil, false
	}
	d.Grant = grant
	if len(jsonData) > 0 {
		d.JSON = string(jsonData)
	}
	return &d, true
}

func (r *PsqlRepo) AddDocument(doc *model.Document) {
	r.db.Exec(context.Background(), "INSERT INTO documents (id, name, mime, file, public, created, grant, owner, json) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)", doc.ID, doc.Name, doc.Mime, doc.File, doc.Public, doc.Created, doc.Grant, doc.Owner, doc.JSON)
}

func (r *PsqlRepo) DeleteDocument(id string) {
	r.db.Exec(context.Background(), "DELETE FROM documents WHERE id=$1", id)
}
