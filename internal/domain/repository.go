package domain

type UserRepo interface {
	Get(login string) (*User, error)
	Add(user *User) error
}

type SessionRepo interface {
	Add(token, login string) error
	Get(token string) (*Session, error)
	Delete(token string) error
}

type DocRepo interface {
	List() ([]*Document, error)
	Get(id string) (*Document, error)
	Add(doc *Document) error
	Delete(id string) error
}
