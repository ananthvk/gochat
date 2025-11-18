package message

import (
	"fmt"
	"net/http"

	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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

func handleGetMessage(m *MessageService, w http.ResponseWriter, r *http.Request) {
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
	message, appErr := m.GetOne(r.Context(), messageId, groupId, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}

	msg := MessageResponse{
		Id:        ulid.ULID(message.ID[:]).String(),
		CreatedAt: message.CreatedAt.Time,
		Type:      message.Type,
		Content:   message.Content,
		GrpId:     ulid.ULID(message.GrpID).String(),
	}

	helpers.RespondWithJSON(w, http.StatusOK, msg)
}

func handleCreateMessage(m *MessageService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot create message without login")
		return
	}
	groupId, err := ulid.Parse(chi.URLParam(r, "group_id"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, "invalid group_id")
		return
	}

	message := MessageCreateRequest{}
	err = helpers.ReadJSONBody(r, &message)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrBadRequest, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(message)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, errs.ErrValidationFailed, fmt.Sprintf("%s", errors))
		return
	}
	id, appErr := m.Create(r.Context(), message.Type, message.Content, groupId, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, map[string]any{"id": id.String()})
}

func handleGetMessages(m *MessageService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot view messages without login")
		return
	}
	groupId, err := ulid.Parse(chi.URLParam(r, "group_id"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.ErrInvalidID, "invalid group_id")
		return
	}
	pagination, err := readPagination(r.URL.Query())
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, errs.ErrValidationFailed, fmt.Sprintf("%s", err))
		return
	}
	msgs, hasMoreBefore, appErr := m.GetAll(r.Context(), pagination, groupId, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	messages := make([]MessageResponse, len(msgs))
	for i, message := range msgs {
		messages[i] = MessageResponse{
			Id:        ulid.ULID(message.ID[:]).String(),
			CreatedAt: message.CreatedAt.Time,
			Type:      message.Type,
			Content:   message.Content,
			GrpId:     ulid.ULID(message.GrpID).String(),
		}
	}
	beforeId := ""
	if len(messages) > 0 {
		beforeId = ulid.ULID(msgs[len(msgs)-1].ID).String()
	}
	helpers.RespondWithJSON(w, 200, map[string]any{
		"messages": messages,
		"cursor": Cursor{
			Before:    beforeId,
			HasBefore: hasMoreBefore,
		},
	})
}

func handleDeleteMessage(m *MessageService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot delete message without login")
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
	appErr := m.Delete(r.Context(), messageId, groupId, user.Id)
	if appErr != nil {
		helpers.RespondWithAppError(w, appErr)
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{"deleted": true})
}
