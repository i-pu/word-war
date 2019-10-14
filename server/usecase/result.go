package usecase

import (
	"errors"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type ResultUsecase interface {
	// userIDはそのuserの結果
	CalResult(userID string) (*entity.Result, error)
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

func (u *resultUsecase) CalResult(userID string) (*entity.Result, error) {
	return nil, errors.New("unimplemented")
}

func (u *resultUsecase) GetResult(userID string) (*entity.Result, error) {
	return nil, errors.New("unimplemented")
}