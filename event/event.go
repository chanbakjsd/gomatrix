package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/matrix"
)

// StrippedEvent represents an event that has been stripped.
// This allows the client to display a room state correctly without its full timeline.
//
// It has the Type, Content, StateKey and Sender field.
type StrippedEvent struct {
	RawEvent
}

// RawEvent represents events that can be sent from homeserver to the client.
type RawEvent struct {
	// Common data for all events.
	Type    Type            `json:"type"`
	Content json.RawMessage `json:"content"`

	// Data that are common for rooms and state events.
	ID               matrix.EventID   `json:"event_id,omitempty"`
	Sender           matrix.UserID    `json:"sender,omitempty"`
	OriginServerTime matrix.Timestamp `json:"origin_server_ts,omitempty"`
	RoomID           matrix.RoomID    `json:"room_id,omitempty"` // NOT included on `/sync` events.
	Unsigned         struct {
		// Age is the time in milliseconds that has elapsed since the event was sent.
		// It is generated by local homeserver and may be incorrect if either server's
		// time is out of sync.
		Age           matrix.Duration `json:"age,omitempty"`
		RedactReason  *RawEvent       `json:"redacted_because,omitempty"`
		TransactionID string          `json:"transaction_id,omitempty"`
	} `json:"unsigned,omitempty"`

	// Data for state events.
	StateKey    string          `json:"state_key,omitempty"`
	PrevContent json.RawMessage `json:"prev_content,omitempty"` // Optional previous content, if available.

	// Data for `m.room.redaction`. The ID of the event that was actually redacted.
	Redacts string `json:"redacts,omitempty"`
}

// Event is a parsed instance of events in Matrix.
type Event interface {
	Type() Type
}

// RoomEvent is an event that is recorded in history and is not one-off.
// Typing is not a RoomEvent for example.
type RoomEvent interface {
	Event

	ID() matrix.EventID
	Room() matrix.RoomID
	Sender() matrix.UserID
	OriginServerTime() matrix.Timestamp
}

// StateEvent is an event that records the change of a state.
type StateEvent interface {
	RoomEvent
	StateKey() string
}

// RoomEventInfo is a helper that satisfies the RoomEvent interface by providing info.
// It does not include Type info.
type RoomEventInfo struct {
	RoomID     matrix.RoomID
	EventID    matrix.EventID
	SenderID   matrix.UserID
	OriginTime matrix.Timestamp
}

// ID satisfies RoomEvent.
func (r RoomEventInfo) ID() matrix.EventID {
	return r.EventID
}

// Room satisfies RoomEvent.
func (r RoomEventInfo) Room() matrix.RoomID {
	return r.RoomID
}

// Sender satisfies RoomEvent.
func (r RoomEventInfo) Sender() matrix.UserID {
	return r.SenderID
}

// OriginServerTime satisfies RoomEvent.
func (r RoomEventInfo) OriginServerTime() matrix.Timestamp {
	return r.OriginTime
}
