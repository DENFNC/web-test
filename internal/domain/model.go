package domain

import "time"

type User struct {
	Login        string
	PasswordHash string
}

type Session struct {
	Token string
	Login string
}

type Document struct {
	ID      string
	Name    string
	Mime    string
	File    bool
	Public  bool
	Created time.Time
	Grant   []string
	Owner   string
	JSON    any
}
