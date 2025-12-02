package message

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/ananthvk/gochat/internal/errs"
	"github.com/ananthvk/gochat/internal/membership"
	"github.com/oklog/ulid/v2"
)

const MessageTypeText = "text"

type MessageService struct {
	Db *database.DatabaseService
}

func NewMessageService(databaseService *database.DatabaseService) *MessageService {
	return &MessageService{
		Db: databaseService,
	}
}

func (m *MessageService) GetOne(ctx context.Context, messageId ulid.ULID, groupId, userId ulid.ULID) (*db.Message, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, m.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(m.Db, ctx, groupId, userId)
	if appErr != nil {
		return nil, appErr
	}

	grp, err := m.Db.Queries.GetMessage(ctx, db.GetMessageParams{ID: messageId[:], GrpID: groupId[:]})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while fetching message", "error", err)
			return nil, errs.Internal("internal server error while fetching message")
		}
		return nil, errs.NotFound("message with the given id not found")
	}
	return grp, nil
}

func (m *MessageService) GetAll(ctx context.Context, pagination Pagination, groupId, userId ulid.ULID) ([]*db.Message, bool, *errs.Error) {
	hasMoreBefore := false
	ctx, cancel := context.WithTimeout(ctx, m.Db.QueryTimeout)
	defer cancel()

	// Note: TOCTOU bug is present, but there's no harm since this is not critical
	// If membership is valid, it means that the group exist, so no need to check for it separately
	appErr := membership.IsUserMemberOfGroup(m.Db, ctx, groupId, userId)
	if appErr != nil {
		return nil, hasMoreBefore, appErr
	}

	var beforeBytes []byte
	if pagination.Before != nil {
		beforeBytes = pagination.Before[:]
	}

	grps, err := m.Db.Queries.GetMessagesInGroup(ctx, db.GetMessagesInGroupParams{
		GrpID:  groupId[:],
		Before: beforeBytes,
		Limit:  int32(pagination.Limit + 1),
	})
	if err != nil {
		slog.Error("internal server error", "errror", err.Error())
		return nil, hasMoreBefore, errs.Internal("internal server error while fetching messages")
	}

	// If we could retrieve limit + 1 rows, it means that hasBefore should be set to true
	if len(grps) == (pagination.Limit + 1) {
		hasMoreBefore = true
		grps = grps[:pagination.Limit]
	}
	return grps, hasMoreBefore, nil
}

func (m *MessageService) Delete(ctx context.Context, messageId, groupId, userId ulid.ULID) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, m.Db.QueryTimeout)
	defer cancel()

	appErr := membership.IsUserMemberOfGroup(m.Db, ctx, groupId, userId)
	if appErr != nil {
		return appErr
	}
	err := m.Db.Queries.DeleteMessage(ctx, db.DeleteMessageParams{ID: messageId[:], GrpID: groupId[:]})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "internal error while deleting message", "error", err)
			return errs.Internal("internal server error while deleting message")
		}
	}
	return nil
}

func (m *MessageService) Create(ctx context.Context, messageType string, content string, groupId, userId ulid.ULID) (ulid.ULID, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, m.Db.QueryTimeout)
	defer cancel()
	id := ulid.Make()

	appErr := membership.IsUserMemberOfGroup(m.Db, ctx, groupId, userId)
	if appErr != nil {
		return ulid.ULID{}, appErr
	}

	_, err := m.Db.Queries.CreateMessage(ctx, db.CreateMessageParams{
		Type:     messageType,
		Content:  content,
		ID:       id[:],
		GrpID:    groupId[:],
		SenderID: userId[:],
	})
	if err != nil {
		slog.ErrorContext(ctx, "internal error while creating message", "error", err)
		return id, errs.Internal("internal server error while creating message")
	}
	return id, nil
}
