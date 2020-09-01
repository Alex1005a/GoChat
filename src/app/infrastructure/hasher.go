package hasher

import (
	"awesomeProject1/src/app/domain"
	"golang.org/x/crypto/bcrypt"
)

type (
	Hasher struct {
	}
)

func NewHasher() domain.Hasher {
	return &Hasher{}
}

func (h *Hasher) PasswordToHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *Hasher) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
