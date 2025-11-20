package app

import (
	"context"
	"time"

	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/config"
	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/group"
	"github.com/ananthvk/gochat/internal/message"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/ananthvk/gochat/internal/token"
)

type App struct {
	Ctx             context.Context
	RealtimeService *realtime.RealtimeService
	DatabaseService *database.DatabaseService
	GroupService    *group.GroupService
	MessageService  *message.MessageService
	AuthService     *auth.AuthService
	TokenService    *token.TokenService
	Config          *config.Config
	Version         string
	StartTime       time.Time
}
