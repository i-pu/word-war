package service

import "github.com/i-pu/word-war/server/domain/repository"

type GameService struct {
	repo repository.GameStateRepository
}

func NewGameService(repo repository.GameStateRepository) *GameService {
	return &GameService{repo: repo}
}