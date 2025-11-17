package group

import (
	"fmt"
	"net/http"

	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/ananthvk/gochat/internal/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

func Routes(g *GroupService, m *message.MessageService) chi.Router {
	router := chi.NewRouter()
	router.Use(auth.Authenticate)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) { handleGetAllGroups(g, w, r) })
	router.Post("/", func(w http.ResponseWriter, r *http.Request) { handleCreateGroup(g, w, r) })
	router.Route("/{group_id}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { handleGetGroup(g, w, r) })
		r.Patch("/", func(w http.ResponseWriter, r *http.Request) { handleUpdateGroup(g, w, r) })
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) { handleDeleteGroup(g, w, r) })
		r.Mount("/message", message.Routes(m))
	})
	return router
}

func handleCreateGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot create group without login")
		return
	}
	grp := GroupCreateRequest{}
	err := helpers.ReadJSONBody(r, &grp)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrBadRequest, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(grp)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, errs.ErrValidationFailed, fmt.Sprintf("%s", errors))
		return
	}

	id, appErr := g.Create(r.Context(), grp.Name, grp.Description, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func handleGetGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot get group without login")
		return
	}
	group_id := chi.URLParam(r, "group_id")
	id, err := ulid.Parse(group_id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, err.Error())
		return
	}
	grp, appErr := g.GetOne(r.Context(), id, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, 200, GroupResponse{
		Id:          ulid.ULID(grp.ID).String(),
		CreatedAt:   grp.CreatedAt,
		Name:        grp.Name,
		Description: grp.Description,
	})
}

func handleDeleteGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot delete group without login")
		return
	}
	group_id := chi.URLParam(r, "group_id")
	id, err := ulid.Parse(group_id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, err.Error())
		return
	}
	appErr := g.Delete(r.Context(), id, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, 200, map[string]any{"deleted": true})
}

func handleUpdateGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot update group without login")
		return
	}
	group_id := chi.URLParam(r, "group_id")
	id, err := ulid.Parse(group_id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, err.Error())
		return
	}
	req := GroupUpdateRequest{}
	err = helpers.ReadJSONBody(r, &req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrBadRequest, err.Error())
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, errs.ErrValidationFailed, fmt.Sprintf("%s", errors))
		return
	}
	grp, appErr := g.Update(r.Context(), req, id, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, 200, GroupResponse{
		Id:          ulid.ULID(grp.ID).String(),
		CreatedAt:   grp.CreatedAt,
		Name:        grp.Name,
		Description: grp.Description,
	})
}

func handleGetAllGroups(g *GroupService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot get all groups without login")
		return
	}
	grps, err := g.GetAll(r.Context(), user.Id)
	if err != nil {
		helpers.RespondWithAppError(w, err)
		return
	}
	groups := make([]GroupResponse, len(grps))
	for i, grp := range grps {
		groups[i] = GroupResponse{
			Id:          ulid.ULID(grp.ID).String(),
			CreatedAt:   grp.CreatedAt,
			Name:        grp.Name,
			Description: grp.Description,
		}
	}
	helpers.RespondWithJSON(w, 200, map[string]any{"groups": groups})
}
