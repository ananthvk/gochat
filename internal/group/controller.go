package group

import (
	"net/http"

	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
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
	createGroupRequest := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}
	err := helpers.ReadJSONBody(r, &createGroupRequest)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "request body malformed", err.Error())
		return
	}

	public_id, err := g.Create(r.Context(), createGroupRequest.Name, createGroupRequest.Description)
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
	// TODO: A more efficient approach will be to create a response struct then marshal into that, do that later

	helpers.RespondWithJSON(w, 200, map[string]any{
		"id":          ulid.ULID(grp.PublicID),
		"created_at":  grp.CreatedAt,
		"name":        grp.Name,
		"description": grp.Description,
	})
}

func handleDeleteGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	g.Delete(r.Context())
	helpers.RespondWithJSON(w, 200, map[string]any{"status": "ok"})
}

func handleUpdateGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	g.Update(r.Context())
	helpers.RespondWithJSON(w, 200, map[string]any{"status": "ok"})
}

func handleGetAllGroups(g *GroupService, w http.ResponseWriter, r *http.Request) {
	g.GetAll(r.Context())
	helpers.RespondWithJSON(w, 200, map[string]any{"status": "ok"})
}
