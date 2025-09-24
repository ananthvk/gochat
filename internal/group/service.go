package group

import (
	"context"
	"log/slog"

	"github.com/ananthvk/gochat/internal/database"
)

type GroupService struct {
	Db *database.DatabaseService
}

func NewGroupService(databaseService *database.DatabaseService) *GroupService {
	return &GroupService{
		Db: databaseService,
	}
}

func (g *GroupService) Create(ctx context.Context) {
	slog.Info("created group")
}

func (g *GroupService) GetOne(ctx context.Context) {
	slog.Info("retrieved a group")
}

func (g *GroupService) Delete(ctx context.Context) {
	slog.Info("deleted a group")
}

func (g *GroupService) Update(ctx context.Context) {
	slog.Info("updated a group")
}

func (g *GroupService) GetAll(ctx context.Context) {
	slog.Info("retrieved all groups for this user")
}
