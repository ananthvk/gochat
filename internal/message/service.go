package message

import "github.com/ananthvk/gochat/internal/database"

type MessageService struct {
	Db database.DatabaseService
}

func NewMessageService(databaseService *database.DatabaseService) *MessageService {
	return &MessageService{
		Db: *databaseService,
	}
}
