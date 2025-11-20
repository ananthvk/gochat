package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/token"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

type AuthService struct {
	Db           *database.DatabaseService
	tokenService *token.TokenService
}

type AuthToken struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

// Token is valid for one day
const DefaultAuthTokenExpiry = 24 * time.Hour

func NewAuthService(databaseService *database.DatabaseService, tokenService *token.TokenService) *AuthService {
	return &AuthService{
		Db:           databaseService,
		tokenService: tokenService,
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

// Returns a login token that can be used to authenticate the user if the password and email matches
func (a *AuthService) LoginByEmail(ctx context.Context, email, plaintextPassword string) (AuthToken, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, a.Db.QueryTimeout)
	defer cancel()

	usr, err := a.Db.Queries.GetUserByEmail(ctx, email)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while login by email", "error", err)
			return AuthToken{}, errs.Internal("internal server error while log in")
		}
		return AuthToken{}, errs.NotFound("user with given email not found")
	}

	p := password{hash: usr.Password}
	ok, err := p.Check(plaintextPassword)
	if err != nil {
		slog.ErrorContext(ctx, "internal error while checking password", "error", err)
		return AuthToken{}, errs.Internal("internal server error while log in")

	}
	if !ok {
		return AuthToken{}, errs.NotAuthenticated("invalid password")
	}

	// Create the authentication token
	token, expiry, appErr := a.tokenService.Create(ctx, token.ScopeAuthenticate, ulid.ULID(usr.ID), DefaultAuthTokenExpiry)
	if appErr != nil {
		return AuthToken{}, appErr
	}
	return AuthToken{Token: token, Expiry: expiry}, nil
}

func (a *AuthService) GetUserById(ctx context.Context, userId ulid.ULID) (*User, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, a.Db.QueryTimeout)
	defer cancel()

	usr, err := a.Db.Queries.GetUserById(ctx, userId[:])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while getting user by id", "error", err)
			return nil, errs.Internal("internal server error while getting user")
		}
		return nil, errs.NotFound("user with given id not found")
	}

	return &User{
		Id:        userId,
		CreatedAt: usr.CreatedAt.Time,
		Name:      usr.Name,
		Username:  usr.Username,
		Email:     usr.Email,
		Password:  password{hash: usr.Password},
		Activated: usr.Activated,
	}, nil
}
