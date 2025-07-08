package model

type User struct {
	Login        string   `json:"login"`
	PasswordHash string   `json:"password_hash"`
	Tokens       []string `json:"tokens,omitempty"`
}
