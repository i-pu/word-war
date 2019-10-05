package repository

type CounterRepository interface {
	IncrCounter() (int64, error)
	SetCounter(value int64) error
	GetCounter() (int64, error)
}
