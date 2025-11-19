package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

type AuthService struct {
	Db *database.DatabaseService
}

func NewAuthService(databaseService *database.DatabaseService) *AuthService {
	return &AuthService{
		Db: databaseService,
	}
}

// Creates a new user, this function also modifies the passed user and sets the ID, CreatedAt, and Password fields if the
// user was created. Set all fields except ID, CreatedAt, and Password. Pass the plaintextPassword, which will get hashed and
// stored in the user
func (a *AuthService) Create(ctx context.Context, user *User, plaintextPassword string) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, a.Db.QueryTimeout)
	defer cancel()
	id := ulid.Make()
	err := user.Password.Set(plaintextPassword)
	if err != nil {
		slog.ErrorContext(ctx, "internal error while creating user", "error", err)
		return errs.Internal("internal server error while creating user")
	}

	row, err := a.Db.Queries.CreateUser(ctx, db.CreateUserParams{
		ID:        id[:],
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Activated: user.Activated,
		Password:  user.Password.hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message)
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "uk_usr_email" {
					return errs.ValidationFailed("a user with the same email already exists")
				} else {
					return errs.ValidationFailed("a user with the same username already exists")
				}
			}
			return errs.Internal("pg error")
		}
		slog.ErrorContext(ctx, "internal error while creating user", "error", err)
		return errs.Internal("internal server error while creating user")
	}
	user.CreatedAt = row.CreatedAt.Time
	user.Id = id
	return nil
}
