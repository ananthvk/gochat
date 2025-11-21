package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/ananthvk/gochat/internal/middleware"
	"github.com/ananthvk/gochat/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func Routes(a *AuthService, middlewares middleware.Middlewares) chi.Router {
	router := chi.NewRouter()
	router.Post("/signup", func(w http.ResponseWriter, r *http.Request) { handleCreateUser(a, w, r) })
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) { handleLoginUserEmail(a, w, r) })
	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticate)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) { handleGetUserInfo(a, w, r) })
		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) { handleLogout(a, w, r) })
	})
	return router
}

func handleCreateUser(a *AuthService, w http.ResponseWriter, r *http.Request) {
	usr := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := helpers.ReadJSONBody(r, &usr); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrBadRequest, err.Error())
		return
	}
	user := &User{
		Name:     usr.Name,
		Email:    usr.Email,
		Username: usr.Username,
		// TODO: Activated is set to true by default, later after implementing email verification, change it to false
		Activated: true,
		Password:  password{Plaintext: &usr.Password},
	}
	v := validator.New(validator.WithRequiredStructEnabled())
	err := v.Struct(user)
	if err != nil {
		helpers.RespondWithAppError(w, errs.ValidationFailed(fmt.Sprintf("%s", err.(validator.ValidationErrors))))
		return
	}
	appErr := a.Create(r.Context(), user, usr.Password)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}

	authToken, expiry, appErr := a.tokenService.Create(r.Context(), token.ScopeAuthenticate, user.Id, DefaultAuthTokenExpiry)

	// Token generation failure is not an error, the user has been created, the client has to try to log in again
	if appErr != nil {
		slog.ErrorContext(r.Context(), "user created, but token generation failed", "error", appErr.String())
		helpers.RespondWithJSON(w, http.StatusCreated, map[string]any{"user": user, "token": nil, "error": "token_generation_failed"})
		return
	}

	// Return a login token after successful signup so that one more round trip need not be made
	resp := struct {
		User  *User     `json:"user"`
		Token AuthToken `json:"token"`
	}{
		User:  user,
		Token: AuthToken{Token: authToken, Expiry: expiry},
	}

	helpers.RespondWithJSON(w, http.StatusCreated, resp)
}

func handleLoginUserEmail(a *AuthService, w http.ResponseWriter, r *http.Request) {
	usr := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=72"`
	}{}
	if err := helpers.ReadJSONBody(r, &usr); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrBadRequest, err.Error())
		return
	}
	v := validator.New(validator.WithRequiredStructEnabled())
	err := v.Struct(usr)
	if err != nil {
		helpers.RespondWithAppError(w, errs.ValidationFailed(fmt.Sprintf("%s", err.(validator.ValidationErrors))))
		return
	}
	token, appErr := a.LoginByEmail(r.Context(), usr.Email, usr.Password)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"authenticate": token})
}

func handleGetUserInfo(a *AuthService, w http.ResponseWriter, r *http.Request) {
	userId, ok := UserIdFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot get profile info without login")
		return
	}
	usr, appErr := a.GetUserById(r.Context(), userId)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"user": usr})
}

func handleLogout(a *AuthService, w http.ResponseWriter, r *http.Request) {
	token, ok := PlaintextTokenFromContext(r.Context())
	if !ok {
		slog.ErrorContext(r.Context(), "error getting token from context")
		helpers.RespondWithAppError(w, errs.Internal("internal server error while logging out"))
		return
	}
	appErr := a.tokenService.DeleteByPlaintextToken(r.Context(), token)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"logout": true})
}
