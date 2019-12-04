package repository

type GameRepository interface {
	UpdateCurrentWord(roomID string, word string) error
	GetCurrentWord(roomID string) (string, error)
}