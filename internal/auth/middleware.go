package auth

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

type ctxKeyUser struct{}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := User{
			Id: ulid.MustParse("01KA9F0Z000000000000000001"),
		}
		ctx := StoreUserInContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(ctxKeyUser{}).(User)
	return u, ok
}

func StoreUserInContext(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, ctxKeyUser{}, u)
}
