package message

import (
	"github.com/ananthvk/gochat/internal/database"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/oklog/ulid/v2"
)

type MessageService struct {
	Db database.DatabaseService
}

func NewMessageService(databaseService *database.DatabaseService) *MessageService {
	return &MessageService{
		Db: *databaseService,
	}
}

func (m *MessageService) GetOne(messageId ulid.ULID, groupId, userId ulid.ULID) (*db.Message, error) {
	return nil, nil
}

func (m *MessageService) GetFromGroup(groupId, userId ulid.ULID) ([]*db.Message, error) {
	return nil, nil
}

func (m *MessageService) Delete(messageId, groupId, userId ulid.ULID) error {
	return nil
}

func (m *MessageService) Create(messageType string, content string, groupId, userId ulid.ULID) (ulid.ULID, error) {
	return ulid.Make(), nil
}
