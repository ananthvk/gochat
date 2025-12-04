package group

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/membership"
	"github.com/jackc/pgx/v5/pgtype"
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

// Create creates a new group, and adds the user who created to the group as a member. If either the creation, or the member addition
// fails, the transaction is rolled back, and an error is returned
func (g *GroupService) Create(ctx context.Context, name, description string, userId ulid.ULID) (ulid.ULID, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	tx, err := g.Db.Pool.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "internal error while creating group", "error", err)
		return ulid.ULID{}, errs.Internal("internal error while creating group")
	}
	defer tx.Rollback(ctx)

	qtx := g.Db.Queries.WithTx(tx)

	id := ulid.Make()

	// Create the group first
	_, err = qtx.CreateGroup(ctx, db.CreateGroupParams{
		Name:        name,
		Description: description,
		ID:          id[:],
		OwnerID:     userId[:],
	})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while creating group", "error", err)
		return ulid.ULID{}, errs.Internal("internal server error while creating group")
	}

	// Add the user who created the group as a member
	_, err = qtx.CreateMembership(ctx, db.CreateMembershipParams{GrpID: id[:], UsrID: userId[:]})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while adding member to group", "error", err)
		return ulid.ULID{}, errs.Internal("internal server error while joining group")
	}

	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "internal error while creating group", "error", err)
		return ulid.ULID{}, errs.Internal("internal server error while creating group")
	}

	return id, nil
}

func (g *GroupService) GetOne(ctx context.Context, groupId, userId ulid.ULID) (*db.Grp, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(g.Db, ctx, groupId, userId)
	if appErr != nil {
		return nil, appErr
	}

	grp, err := g.Db.Queries.GetGroup(ctx, groupId[:])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while fetching group", "error", err)
			return nil, errs.Internal("internal server error while fetching group")
		}
		return nil, errs.NotFound("group with the given id not found")
	}
	return grp, nil
}

func (g *GroupService) Delete(ctx context.Context, groupId, userId ulid.ULID) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(g.Db, ctx, groupId, userId)
	if appErr != nil {
		return appErr
	}

	err := g.Db.Queries.DeleteGroup(ctx, groupId[:])
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while deleting group", "error", err)
			return errs.Internal("internal server error while deleting group")
		}
	}
	return nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Updates a group (supports partial updates)
func (g *GroupService) Update(ctx context.Context, req GroupUpdateRequest, groupId, userId ulid.ULID) (*db.Grp, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(g.Db, ctx, groupId, userId)
	if appErr != nil {
		return nil, appErr
	}

	group, err := g.Db.Queries.UpdateGroupById(ctx, db.UpdateGroupByIdParams{
		Name:        pgtype.Text{String: deref(req.Name), Valid: req.Name != nil},
		Description: pgtype.Text{String: deref(req.Description), Valid: req.Description != nil},
		ID:          groupId[:],
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while updating group", "error", err)
			return nil, errs.Internal("internal server error while updating group")
		}
		return nil, errs.NotFound("group with given id not found")
	}
	return group, nil
}

func (g *GroupService) GetAll(ctx context.Context, userId ulid.ULID) ([]*db.GetGroupsRow, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	grps, err := g.Db.Queries.GetGroups(ctx, userId[:])
	if err != nil {
		return nil, errs.Internal("internal server error while fetching groups")
	}
	return grps, nil
}

func (g *GroupService) AddMemberToGroup(ctx context.Context, groupId, userId ulid.ULID) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(g.Db, ctx, groupId, userId)
	if appErr == nil {
		// If the user already exists, don't do anything, idempotent behavior
		return nil
	}

	// Only return the error if it's an internal error
	if appErr.Kind == errs.ErrInternal {
		return appErr
	}

	_, err := g.Db.Queries.CreateMembership(ctx, db.CreateMembershipParams{GrpID: groupId[:], UsrID: userId[:]})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while adding member to group", "error", err)
		return errs.Internal("internal server error while joining group")
	}
	return nil
}

func (g *GroupService) GetMembers(ctx context.Context, groupId, userId ulid.ULID) ([]*db.GetGroupMembersWithNameRow, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, g.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(g.Db, ctx, groupId, userId)
	if appErr != nil {
		return nil, appErr
	}
	members, err := g.Db.Queries.GetGroupMembersWithName(ctx, groupId[:])
	if err != nil {
		slog.ErrorContext(ctx, "internal error while fetching members of the group", "error", err)
		return nil, errs.Internal("internal server error while fetching members")
	}
	return members, nil
}
