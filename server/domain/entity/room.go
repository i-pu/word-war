package entity

type Room struct {
	RoomID         string `json:"roomID" validate:"required"`
	CurrentMessage *Message `json:"currentMessage" validate: "required"`
}
