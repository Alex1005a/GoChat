package domain

type User struct {
	Id      string              `bson:"_id" json:"id,omitempty"`
	Name   string               `json:"name"`
	PasswordHash  string        `json:"hash"`
}

type UserRepo interface {
	CreateUser(name string, password string) (string, error)

	UserByID(id string) (User, error)

	UserByUsername(name string) (User, error)

	DeleteUserByID(id string) error
}

type Hasher interface {
	PasswordToHash(password string) (string, error)

	CheckPassword(password, hash string) bool
}