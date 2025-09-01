package realtime

type Event interface {
}

type RegisterConnectionEvent struct {
	Client *Client
}

type UnregisterConnectionEvent struct {
	Client *Client
}

type DataEvent struct {
	Client *Client
	Data   []byte
}
