package service

import "github.com/i-pu/word-war/server/domain/repository"

type GameService struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) *GameService {
	return &GameService{repo: repo}
}