package usecase

import (
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type CounterUsecase interface {
	Init(counter *entity.Counter) (*entity.Counter, error)
	Incr() (*entity.Counter, error)
	Get() (*entity.Counter, error)
}

type counterUsecase struct {
	repo    repository.CounterRepository
	service *service.CounterService
}

func NewCounterUsecase(repo repository.CounterRepository, service *service.CounterService) *counterUsecase {
	return &counterUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *counterUsecase) Init(counter *entity.Counter) (*entity.Counter, error) {
	if err := u.repo.SetCounter(counter.Value); err != nil {
		return nil, err
	}
	return counter, nil
}

func (u *counterUsecase) Incr() (*entity.Counter, error) {
	value, err := u.repo.IncrCounter()
	if err != nil {
		return nil, err
	}
	return &entity.Counter{Value: value}, nil
}
func (u *counterUsecase) Get() (*entity.Counter, error) {
	value, err := u.repo.GetCounter()
	if err != nil {
		return nil, err
	}
	return &entity.Counter{Value: value}, nil
}
