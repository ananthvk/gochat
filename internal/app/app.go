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

func NewApp(ctx context.Context, cfg *config.Config, version string) (*App, error) {
	dbService, err := database.NewDatabaseService(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tokenService := token.NewTokenService(dbService)
	authService := auth.NewAuthService(dbService, tokenService)
	groupService := group.NewGroupService(dbService)
	messageService := message.NewMessageService(dbService)
	realtimeService := realtime.NewRealtimeService(ctx)

	app := &App{
		Ctx:             ctx,
		RealtimeService: realtimeService,
		DatabaseService: dbService,
		GroupService:    groupService,
		MessageService:  messageService,
		AuthService:     authService,
		TokenService:    tokenService,
		Config:          cfg,
		Version:         version,
		StartTime:       time.Now(),
	}

	return app, nil
}
