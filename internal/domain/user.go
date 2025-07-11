package domain

import "time"

type User struct {
	ID        string
	Login     string
	Password  string
	CreatedAt time.Time
}

type UserCredentials struct {
	Password string
	Token    string
}
