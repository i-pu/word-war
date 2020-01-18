package entity

type GameState struct {
	RoomID      string `json:"roomID" validate:"required"`
	CurrentWord string `json:"currentWord" validate: "required"`
}
