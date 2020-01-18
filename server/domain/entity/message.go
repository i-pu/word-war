package entity

type Message struct {
	Message string `json:"message" validate:"required"`
	UserID  string `json:"userID" validate:"required"`
	RoomID  string `json:"roomID" validate:"required"`
}
