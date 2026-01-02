package middleware

import "net/http"

type Middlewares struct {
	Authenticate           func(http.Handler) http.Handler
	AuthenticateQueryParam func(http.Handler) http.Handler
}
