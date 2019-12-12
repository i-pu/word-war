package entity

type GameState struct {
	RoomID      string `json:"roomID" validate:"required"`
	CurrentWord string `json:"currentWord" validate: "required"`
	// TODO
	// history []Message `json:"history"`
	// member []User `json: "member"`
	// created Time `json: "created"`
}
