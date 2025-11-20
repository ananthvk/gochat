package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
)

const defaultTokenSize = 24 // Size of the token in bytes
const ScopeAuthenticate = "authenticate"

type TokenService struct {
	Db *database.DatabaseService
}

func NewTokenService(databaseService *database.DatabaseService) *TokenService {
	return &TokenService{
		Db: databaseService,
	}
}

// This function returns the plaintext token and the expiry time
func (t *TokenService) Create(ctx context.Context, scope string, userId ulid.ULID, ttl time.Duration) (string, time.Time, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, t.Db.QueryTimeout)
	defer cancel()

	randomToken := make([]byte, defaultTokenSize)
	_, err := rand.Read(randomToken)
	if err != nil {
		// Shouldn't fail
		slog.ErrorContext(ctx, "internal error while creating token", "error", err)
		return "", time.Time{}, errs.Internal("internal server error while creating token")
	}

	plaintextToken := base64.RawURLEncoding.EncodeToString(randomToken[:])
	hashedToken := sha256.Sum256(randomToken)
	expiry := time.Now().Add(ttl)

	err = t.Db.Queries.CreateToken(ctx, db.CreateTokenParams{
		UsrID:  userId[:],
		Scope:  scope,
		Hash:   hashedToken[:],
		Expiry: pgtype.Timestamptz{Time: expiry, Valid: true},
	})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while creating token", "error", err)
		return "", time.Time{}, errs.Internal("internal server error while creating token")
	}

	return plaintextToken, expiry, nil
}

// Verify checks if the token is valid, it returns an error only if the token is invalid or if an internal error occured
// If the token exists, a user id is returned, otherwise ulid.ULID{} is returned
func (t *TokenService) Verify(ctx context.Context, scope, plaintextToken string) (ulid.ULID, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, t.Db.QueryTimeout)
	defer cancel()

	plaintext, err := base64.RawURLEncoding.DecodeString(plaintextToken)
	if err != nil {
		return ulid.ULID{}, errs.BadRequest("invalid token format")
	}
	hashedToken := sha256.Sum256(plaintext)

	id, err := t.Db.Queries.VerifyToken(ctx, db.VerifyTokenParams{
		Scope: scope,
		Hash:  hashedToken[:],
	})

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while verifying token", "error", err)
			return ulid.ULID{}, errs.Internal("internal server error verifying token")
		}
		return ulid.ULID{}, nil
	}

	return ulid.ULID(id), nil
}

func (t *TokenService) DeleteByHash(ctx context.Context, hash []byte) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, t.Db.QueryTimeout)
	defer cancel()

	err := t.Db.Queries.DeleteTokenByHash(ctx, hash)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while deleting token by hash", "error", err)
			return errs.Internal("internal server error while deleting token")
		}
		// If the token is already deleted, it's not an error
	}
	return nil
}

func (t *TokenService) DeleteByPlaintextToken(ctx context.Context, plaintextToken string) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, t.Db.QueryTimeout)
	defer cancel()

	plaintext, err := base64.RawURLEncoding.DecodeString(plaintextToken)
	if err != nil {
		return errs.BadRequest("invalid token format")
	}
	hashedToken := sha256.Sum256(plaintext)

	err = t.Db.Queries.DeleteTokenByHash(ctx, hashedToken[:])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while deleting token by hash", "error", err)
			return errs.Internal("internal server error while deleting token")
		}
		// If the token is already deleted, it's not an error
	}
	return nil
}

func (t *TokenService) DeleteAllByUserId(ctx context.Context, scope string, userId ulid.ULID) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, t.Db.QueryTimeout)
	defer cancel()
	err := t.Db.Queries.DeleteTokensByUserId(ctx, db.DeleteTokensByUserIdParams{UsrID: userId[:], Scope: scope})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while deleting all tokens for user", "error", err)
			return errs.Internal("internal server error while deleting token")
		}
	}
	return nil
}
