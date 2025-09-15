package realtime

import "github.com/google/uuid"

type room struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type RealtimeService struct {
	rooms map[string]*room
	hub   *hub
}

func NewRealtimeService() *RealtimeService {
	return &RealtimeService{
		rooms: map[string]*room{},
		hub:   newHub(),
	}
}

// CreateRoom creates a new room, and returns it. If a room with the same name already exists, it returns
// the existing room.
func (r *RealtimeService) CreateRoom(name string) *room {
	existingRoom, ok := r.rooms[name]
	if !ok {
		room := &room{Id: uuid.New(), Name: name}
		r.rooms[name] = room
		return room
	}
	return existingRoom
}

// ListRooms returns a list of all the rooms on this server. If there are no rooms, an empty slice is returned
func (r *RealtimeService) ListRooms() []*room {
	rooms := make([]*room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// GetRoomByName returns the room with the given name. Note: In the current implementation, rooms have unique names.
// If the room with the given name is not found, nil is returned
func (r *RealtimeService) GetRoomByName(name string) *room {
	room := r.rooms[name]
	return room
}

// StartHubEventLoop starts the event loop of the hub in a separate goroutine
func (r *RealtimeService) StartHubEventLoop() {
	go r.hub.RunEventLoop()
}
