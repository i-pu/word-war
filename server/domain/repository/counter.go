package repository

type CounterRepository interface {
	IncrCounter(roomID string) (int64, error)
	SetCounter(roomID string, value int64) error
	GetCounter(roomID string) (int64, error)
}
