package group

import (
	"fmt"
	"net/http"

	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

func Routes(g *GroupService) chi.Router {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) { handleGetAllGroups(g, w, r) })
	router.Post("/", func(w http.ResponseWriter, r *http.Request) { handleCreateGroup(g, w, r) })
	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) { handleGetGroup(g, w, r) })
	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { handleDeleteGroup(g, w, r) })
	router.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) { handleUpdateGroup(g, w, r) })
	return router
}

func handleCreateGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	grp := GroupCreateRequest{}
	err := helpers.ReadJSONBody(r, &grp)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "request body malformed", err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(grp)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, "validation failed", fmt.Sprintf("%s", errors))
		return
	}

	public_id, err := g.Create(r.Context(), grp.Name, grp.Description)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "group not created", err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, map[string]any{"id": public_id})
}

func handleGetGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	public_id := chi.URLParam(r, "id")
	id, err := ulid.Parse(public_id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	grp, err := g.GetOne(r.Context(), id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "group not found", err.Error())
		return
	}
	helpers.RespondWithJSON(w, 200, GroupResponse{
		Id:          ulid.ULID(grp.PublicID).String(),
		CreatedAt:   grp.CreatedAt,
		Name:        grp.Name,
		Description: grp.Description,
	})
}

func handleDeleteGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	public_id := chi.URLParam(r, "id")
	id, err := ulid.Parse(public_id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	err = g.Delete(r.Context(), id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "delete failed", err.Error())
		return
	}
	helpers.RespondWithJSON(w, 200, map[string]any{"deleted": true})
}

func handleUpdateGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	public_id := chi.URLParam(r, "id")
	id, err := ulid.Parse(public_id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid id", err.Error())
		return
	}
	req := GroupUpdateRequest{}
	err = helpers.ReadJSONBody(r, &req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "request body malformed", err.Error())
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, "validation failed", fmt.Sprintf("%s", errors))
		return
	}
	grp, err := g.Update(r.Context(), id, req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "error while updating group", err.Error())
		return
	}
	helpers.RespondWithJSON(w, 200, GroupResponse{
		Id:          ulid.ULID(grp.PublicID).String(),
		CreatedAt:   grp.CreatedAt,
		Name:        grp.Name,
		Description: grp.Description,
	})
}

func handleGetAllGroups(g *GroupService, w http.ResponseWriter, r *http.Request) {
	grps, err := g.GetAll(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "error getting records", err.Error())
		return
	}
	groups := make([]GroupResponse, len(grps))
	for i, grp := range grps {
		groups[i] = GroupResponse{
			Id:          ulid.ULID(grp.PublicID).String(),
			CreatedAt:   grp.CreatedAt,
			Name:        grp.Name,
			Description: grp.Description,
		}
	}
	helpers.RespondWithJSON(w, 200, map[string]any{"groups": groups})
}
