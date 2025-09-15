package realtime

import "github.com/google/uuid"

type hubRoom struct {
	Id      uuid.UUID
	Clients map[uuid.UUID]struct{}
}

func newHubRoom(id uuid.UUID) *hubRoom {
	return &hubRoom{
		Id:      id,
		Clients: map[uuid.UUID]struct{}{},
	}
}
