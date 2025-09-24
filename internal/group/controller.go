package group

import (
	"net/http"

	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
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
	g.Create(r.Context())
	helpers.RespondWithJSON(w, 200, map[string]any{"status": "ok"})
}

func handleGetGroup(g *GroupService, w http.ResponseWriter, r *http.Request) {
	g.GetOne(r.Context())
	helpers.RespondWithJSON(w, 200, map[string]any{"status": "ok"})
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
