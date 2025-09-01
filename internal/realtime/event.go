package realtime

type event any

type registerConnectionEvent struct {
	Client *client
}

type unregisterConnectionEvent struct {
	Client *client
}

type dataEvent struct {
	Client *client
	Data   []byte
}
