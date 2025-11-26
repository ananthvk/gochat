package membership

import (
	"context"
	"log/slog"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/oklog/ulid/v2"
)

// Returns an error if the user is not a member of the group or if there was an error fetching the status
// Otherwise returns nil
func IsUserMemberOfGroup(databaseService *database.DatabaseService, ctx context.Context, groupId ulid.ULID, userId ulid.ULID) *errs.Error {
	isMember, err := databaseService.Queries.CheckMembership(ctx, db.CheckMembershipParams{GrpID: groupId[:], UsrID: userId[:]})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while checking membership", "error", err)
		return errs.Internal("internal server error while fetching group")
	}
	if !isMember {
		return errs.NotAuthorized("not authorized to view details of the group")
	}
	return nil
}
