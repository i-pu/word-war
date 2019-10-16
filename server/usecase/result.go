package usecase

import (
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type ResultUsecase interface {
	// userIDはそのuserの結果
	IncrResult(userID string, by int64) error
	GetResult(userID string) (*entity.Result, error)
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

func (u *resultUsecase) IncrResult(userID string, by int64) error {
	return u.repo.IncrBy(userID, by)
}

func (u *resultUsecase) GetResult(userID string) (*entity.Result, error) {
	return u.repo.Get(userID)
}
