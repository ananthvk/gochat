package auth

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

const hashCost = 12

type User struct {
	Id        ulid.ULID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
	Username  string    `json:"username" validate:"required,min=3,max=32"`
	Email     string    `json:"email" validate:"required,email"`
	Password  password  `json:"-" validate:"required"`
	Activated bool      `json:"activated"`
}

type password struct {
	Plaintext *string `json:"-" validate:"required,min=8,max=72"`
	hash      []byte
}

func (p *password) Set(passwordPlaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordPlaintext), hashCost)
	if err != nil {
		return err
	}
	p.Plaintext = &passwordPlaintext
	p.hash = hash
	return nil
}

// Checks if the plaintext password matches the stored hashed password
// Returns true if they match, false otherwise
func (p *password) Check(passwordPlaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(passwordPlaintext))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u *User) Validate() error {
	validator := validator.New(validator.WithRequiredStructEnabled())
	if u.Password.hash == nil {
		panic("logic error: hash is nil")
	}
	return validator.Struct(u)
}
