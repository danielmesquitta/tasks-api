package hasher

import (
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) Hash(plaintext string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), 14)
	if err != nil {
		return "", entity.NewErr(err)
	}
	return string(bytes), nil
}

func (b *Bcrypt) Match(plaintext, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}
