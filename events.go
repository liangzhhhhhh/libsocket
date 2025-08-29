package libsocket

type (
	EventType int
)

const (
	EventConnect EventType = iota
	EventReconnect
	EventClose
)
