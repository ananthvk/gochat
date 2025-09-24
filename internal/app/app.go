package app

import (
	"context"
	"time"

	"github.com/ananthvk/gochat/internal/config"
	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/realtime"
)

type App struct {
	Ctx             context.Context
	RealtimeService *realtime.RealtimeService
	DatabaseService *database.DatabaseService
	Config          *config.Config
	Version         string
	StartTime       time.Time
}
