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

type ClientMessage struct {
	Message string `json:"message" validate:"required"`
	UserID  string `json:"userID" validate:"required"`
	RoomID  string `json:"roomID" validate:"required"`
}

type ServerMessage struct {
	Done bool
}
