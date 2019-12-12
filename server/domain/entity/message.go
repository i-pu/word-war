package entity

type Message struct {
	Message string `json:"message" validate:"required"`
	UserID  string `json:"userID" validate:"required"`
	RoomID  string `json:"roomID" validate:"required"`

	// TODO: サーバ
	// Type    string `json:"type" validate:"required"`
	// ClientMessage *ClientMessage
	// ServerMessage *ServerMessage
}

// TODO: もしかしたら必要になるかも
type ClientMessage struct {
	Message string `json:"message" validate:"required"`
	UserID  string `json:"userID" validate:"required"`
	RoomID  string `json:"roomID" validate:"required"`
}

// TODO: もしかしたら必要になるかも
type ServerMessage struct {
	Done bool
}
