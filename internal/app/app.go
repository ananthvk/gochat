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
)

type App struct {
	Ctx             context.Context
	RealtimeService *realtime.RealtimeService
	DatabaseService *database.DatabaseService
	GroupService    *group.GroupService
	MessageService  *message.MessageService
	AuthService     *auth.AuthService
	Config          *config.Config
	Version         string
	StartTime       time.Time
}
