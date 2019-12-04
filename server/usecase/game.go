package usecase

import (
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type GameUsecase interface {
	Init(roomID string, game *entity.Game) (*entity.Game, error)
	UpdateCurrentWord(roomID string, word string) (*entity.Game, error)
}

type gameUsecase struct {
	repo	repository.GameRepository
	service *service.GameService
}

func NewGameUsecase(repo repository.GameRepository, service *service.GameService) *gameUsecase {
	return &gameUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *gameUsecase) Init(roomID string, game *entity.Game) (*entity.Game, error) {
	err := u.repo.UpdateCurrentWord(roomID, "")
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (u *gameUsecase) UpdateCurrentWord(roomID string, word string) (*entity.Game, error) {
	err := u.repo.UpdateCurrentWord(roomID, word)
	if err != nil {
		return nil, err
	}
	return &entity.Game{RoomID:roomID,CurrentWord: word}, nil
}