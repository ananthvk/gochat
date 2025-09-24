package group

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/oklog/ulid/v2"
)

type GroupService struct {
	Db *database.DatabaseService
}

func NewGroupService(databaseService *database.DatabaseService) *GroupService {
	return &GroupService{
		Db: databaseService,
	}
}

func (g *GroupService) Create(ctx context.Context, name, description string) (ulid.ULID, error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	public_id := ulid.Make()

	_, err := g.Db.Queries.CreateGroup(ctx, db.CreateGroupParams{
		Name:        name,
		Description: description,
		PublicID:    public_id[:],
	})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while creating group", "error", err)
		return public_id, err
	}
	return public_id, nil
}

func (g *GroupService) GetOne(ctx context.Context, public_id ulid.ULID) (*db.Grp, error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()
	grp, err := g.Db.Queries.GetGroupByPublicId(ctx, public_id[:])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while fetching group", "error", err)
			return nil, errors.New("internal server error")
		}
		return nil, errors.New("group with the given id not found")
	}
	return grp, nil
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
