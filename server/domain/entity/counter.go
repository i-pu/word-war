package entity

type Counter struct {
	RoomID string `json:"roomID" validate:"required"`
	Value  int64
}
