package middleware

import "net/http"

type Middlewares struct {
	Authenticate func(http.Handler) http.Handler
}
