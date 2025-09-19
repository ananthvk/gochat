package app

import (
	"context"

	"github.com/ananthvk/gochat/internal/realtime"
)

type App struct {
	Ctx             context.Context
	RealtimeService *realtime.RealtimeService
}
