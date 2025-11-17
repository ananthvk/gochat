package message

import (
	"net/http"

	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func Routes(m *MessageService) chi.Router {
	router := chi.NewRouter()
	router.Use(auth.Authenticate)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) { handleGetMessages(m, w, r) })
	router.Post("/", func(w http.ResponseWriter, r *http.Request) { handleCreateMessage(m, w, r) })
	router.Route("/{message_id}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { handleGetMessage(m, w, r) })
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) { handleDeleteMessage(m, w, r) })
	})
	return router
}

func handleGetMessages(m *MessageService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot get message without login")
		return
	}
	groupId, err := ulid.Parse(chi.URLParam(r, "group_id"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, "invalid group_id")
		return
	}
	messageId, err := ulid.Parse(chi.URLParam(r, "message_id"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, "invalid message_id")
		return
	}

	_, err = m.GetOne(messageId, groupId, user.Id)

	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"status": "ok", "group_id": groupId, "message_id": messageId})
}

func handleCreateMessage(_ *MessageService, w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "group_id")
	messageId := chi.URLParam(r, "message_id")
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"status": "ok", "group_id": groupId, "message_id": messageId})
}

func handleGetMessage(m *MessageService, w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "group_id")
	messageId := chi.URLParam(r, "message_id")
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"status": "ok", "group_id": groupId, "message_id": messageId})
}

func handleDeleteMessage(m *MessageService, w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "group_id")
	messageId := chi.URLParam(r, "message_id")
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"status": "ok", "group_id": groupId, "message_id": messageId})
}
