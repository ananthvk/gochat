package auth

import (
	"fmt"
	"net/http"

	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func Routes(a *AuthService) chi.Router {
	router := chi.NewRouter()
	router.Post("/signup", func(w http.ResponseWriter, r *http.Request) { handleCreateUser(a, w, r) })
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
	helpers.RespondWithJSON(w, 200, map[string]any{"user": user})
}
