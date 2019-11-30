package usecase

import (
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type ResultUsecase interface {
	// userIDはそのuserの結果
	IncrResult(roomID string, userID string, by int64) error
	GetResult(roomID string, userID string) (*entity.Result, error)
}

type resultUsecase struct {
	repo    repository.ResultRepository
	service *service.ResultService
}

func NewResultUsecase(repo repository.ResultRepository, service *service.ResultService) *resultUsecase {
	return &resultUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *resultUsecase) IncrResult(roomID string, userID string, by int64) error {
	return u.repo.IncrBy(roomID, userID, by)
}

func (u *resultUsecase) GetResult(roomID string, userID string) (*entity.Result, error) {
	return u.repo.Get(roomID, userID)
}
