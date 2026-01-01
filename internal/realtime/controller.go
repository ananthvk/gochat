package realtime

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/ananthvk/gochat/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Routes(rt *RealtimeService, middlewares middleware.Middlewares) chi.Router {
	realtimeRouter := chi.NewRouter()
	realtimeRouter.Use(middlewares.Authenticate)
	realtimeRouter.Get("/ws", func(w http.ResponseWriter, r *http.Request) { handlerCreateWSConnection(rt, w, r) })
	return realtimeRouter
}

func handlerCreateWSConnection(rt *RealtimeService, w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.UserIdFromContext(r.Context())
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, errs.ErrNotAuthenticated, "cannot connect to ws notifications without login")
		return
	}

	// TODO: Fix this to check origin correctly, also install and use cors package
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.ErrorContext(r.Context(), "websocket upgrade failed", "error", err)
		helpers.RespondWithError(w, http.StatusUpgradeRequired, "websocket upgrade failed", err.Error())
		return
	}

	clientId := rt.RegisterConnection(conn, ulid.Make())

	// TOOD: Move it into a service later, for now just do the db call here to get the list of groups the user is part of
	// TODO: Also move it into a separate notification service (which depends on db + realtime instead of making realtime depend on the db)
	ctx, cancel := context.WithTimeout(r.Context(), rt.Db.QueryTimeout)
	defer cancel()

	grps, err := rt.Db.Queries.GetGroups(ctx, userId[:])
	if err != nil {
		helpers.RespondWithAppError(w, errs.Internal("internal server error while fetching groups"))
		return
	}
	groupIds := make([]ulid.ULID, len(grps))
	for i, grp := range grps {
		groupIds[i] = ulid.ULID(grp.ID)
	}
	rt.AddConnectionToRooms(groupIds, clientId)
}
