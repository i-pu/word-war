package repository

type GameStateRepository interface {
	InitWord(roomID string, word string) error
	LockCurrentWord(roomID string) error
	UnlockCurrentWord(roomID string) error
	UpdateCurrentWord(roomID string, word string) error
	GetCurrentWord(roomID string) (string, error)
	AddUser(roomID string, userID string) error
	GetUsers(roomID string) ([]string, error)
}
