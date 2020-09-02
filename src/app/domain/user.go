package domain

import "github.com/dgrijalva/jwt-go"

type Token struct {
	UserId string
	jwt.StandardClaims
}

type User struct {
	Id           string `bson:"_id" json:"id,omitempty"`
	Name         string `json:"name"`
	PasswordHash string `json:"hash"`
	Token        string `json:"token"`
}

type UserRepo interface {
	CreateUser(name string, password string) (string, error)

	UserByID(id string) (User, error)

	UserByUsername(name string) (User, error)

	DeleteUserByID(id string) error

	Login(name, password string) User
}

type Hasher interface {
	PasswordToHash(password string) (string, error)

	CheckPassword(password, hash string) bool
}
