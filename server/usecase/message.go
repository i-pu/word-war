package usecase

import (
	"context"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
)

type MessageUsecase interface {
	JudgeMessage(message *entity.Message) bool
	SendMessage(message *entity.Message) error
	GetMessage(ctx context.Context) (<-chan *entity.Message, <-chan error)
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

// test usecase
func (u *messageUsecase) JudgeMessage(message *entity.Message) bool {
	return u.repo.IsSingleNoun(message)
}

func (u *messageUsecase) SendMessage(message *entity.Message) error {
	// [TODO] validation message

	if err := u.repo.Publish(message); err != nil {
		return err
	}
	return nil
}

// GetMessage ctx is used to get cancel signal from parent to cancel pub/sub job
// , so this ctx must be child context.
func (u *messageUsecase) GetMessage(ctx context.Context) (<-chan *entity.Message, <-chan error) {
	messageChan, errChan := u.repo.Subscribe(ctx)
	return messageChan, errChan
}
