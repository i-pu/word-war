package service

import "github.com/i-pu/word-war/server/domain/repository"

type CounterService struct {
	repo repository.CounterRepository
}

func NewCounterService(repo repository.CounterRepository) *CounterService {
	return &CounterService{repo: repo}
}
