package repository

type RoomRepository interface {
	Lock() error
	Unlock() error
	GetRoomCandidates() ([]string, error)
	AddRoomCandidate(roomID string) error
	DeleteRoomCandidate(roomID string) error
	DeleteRoom(roomID string) error
}