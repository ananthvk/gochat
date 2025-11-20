package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/ananthvk/gochat/internal/token"
	"github.com/oklog/ulid/v2"
)

const HardcodedUserId = "01KA9F0Z000000000000000001"

type ctxKeyUser struct{}

type ctxKeyToken struct{}

// AuthMiddleware is a factory that creates an auth middleware by wrapping the token service through a closure
func AuthMiddleware(tokenService *token.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// To inform the browser that the cached value depends upon the value of the Authorization header
			w.Header().Add("Vary", "Authorization")

			// Check if the authorization header is present
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				helpers.RespondWithAppError(w, errs.NotAuthenticated("authorization header is not present"))
				return
			}

			// If the authorization header is present, parse it in the format Bearer [token], if it's not in this format,
			// return a 401 response
			authorizationParts := strings.Split(authorization, " ")
			if len(authorizationParts) != 2 {
				helpers.RespondWithAppError(w, errs.NotAuthenticated("authorization header is malformed"))
				return
			}
			if authorizationParts[0] != "Bearer" {
				helpers.RespondWithAppError(w, errs.NotAuthenticated("authorization header is malformed, the first token is not 'Bearer'"))
				return
			}

			userId, appErr := tokenService.Verify(r.Context(), token.ScopeAuthenticate, authorizationParts[1])
			if appErr != nil {
				appErr.Kind = errs.ErrNotAuthenticated
				appErr.Status = http.StatusUnauthorized
				helpers.RespondWithAppError(w, appErr)
				return
			}

			if userId == (ulid.ULID{}) {
				helpers.RespondWithAppError(w, errs.NotAuthenticated("invalid or expired authentication token"))
				return
			}

			ctx := StoreUserIdInContext(r.Context(), userId)
			ctx = StorePlaintextTokenInContext(ctx, authorizationParts[1])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIdFromContext(ctx context.Context) (ulid.ULID, bool) {
	u, ok := ctx.Value(ctxKeyUser{}).(ulid.ULID)
	return u, ok
}

func PlaintextTokenFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(ctxKeyToken{}).(string)
	return u, ok
}

func StoreUserIdInContext(ctx context.Context, userId ulid.ULID) context.Context {
	return context.WithValue(ctx, ctxKeyUser{}, userId)
}

func StorePlaintextTokenInContext(ctx context.Context, plaintextToken string) context.Context {
	return context.WithValue(ctx, ctxKeyToken{}, plaintextToken)
}
