package models

import "github.com/gorilla/websocket"

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
	EVENT_CHANGE
)

type Subscriber struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      string
	ToUser    string
	Timestamp int // Unix timestamp (secs)
	Content   string
	Message   *Message
}

type Message struct {
	FromChatRoom string `json:"fromChatRoom"`
	Text         string `json:"text"`
}

type ReceivedMessage struct {
	FromName string `json:"fromName"`
	ToName   string `json:"toName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
