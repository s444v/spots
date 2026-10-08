package signaling

import "encoding/json"

type wsMessage struct {
	Type    MessageType     `json:"type"`
	Room    string          `json:"room,omitempty"`
	Name    string          `json:"name,omitempty"`
	ID      string          `json:"id,omitempty"`
	To      string          `json:"to,omitempty"`
	From    string          `json:"from,omitempty"`
	Peers   []Peer          `json:"peers,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

type Peer struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type MessageType string

const (
	TypeJoin       MessageType = "join"
	TypeLeave      MessageType = "leave"
	TypeSignal     MessageType = "signal"
	TypeWelcome    MessageType = "welcome"
	TypePeerJoined MessageType = "peer-joined"
	TypePeerLeft   MessageType = "peer-left"
	TypeError      MessageType = "error"
)
