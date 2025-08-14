package chat

import "encoding/json"

type Message struct {
	Type    string          `json:"type"`              // e.g. "chat","system","typing"
	From    string          `json:"from,omitempty"`    // sender ID
	To      string          `json:"to,omitempty"`      // optional direct recipient
	Room    string          `json:"room,omitempty"`    // optional room/topic
	Payload json.RawMessage `json:"payload,omitempty"` // defer decoding

}
