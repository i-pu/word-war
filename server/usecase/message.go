package usecase

import (
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type MessageUsecase interface {
	SendMessage(key string, message *entity.Message) error
	GetMessage(key string) (*entity.Message, error)
}

type messageUsecase struct {
	repo    repository.MessageRepository
	service *service.MessageService
}

func NewMessageUsecase(repo repository.MessageRepository, service *service.MessageService) *messageUsecase {
	return &messageUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *messageUsecase) SendMessage(key string, message *entity.Message) error {
	if err := u.repo.Publish("message", message); err != nil {
		return err
	}
	return nil
}

func (u *messageUsecase) GetMessage(key string) (*entity.Message, error) {
	return nil, nil
}

func (u *messageUsecase) GetNowCounter() (int64, error) {
	return 3, nil
}
