package service

import (
	"github.com/i-pu/word-war/server/domain/repository"
)

type ResultService struct {
	repo repository.ResultRepository
}

func NewResultService(repo repository.ResultRepository) *ResultService {
	return &ResultService{repo: repo}
}
