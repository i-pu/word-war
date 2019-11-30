package usecase

import (
	"context"
	"errors"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
	"regexp"
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

// ひらがな && 1単語 && 名詞
func (u *messageUsecase) JudgeMessage(message *entity.Message) bool {
	// TODO: 今なんの単語からはじまるのか
	// TODO: `ん`で終わるかどうか
	r := regexp.MustCompile(`^\p{Hiragana}+$`)
	return r.Match([]byte(message.Message)) && u.repo.IsSingleNoun(message)
}

func (u *messageUsecase) SendMessage(message *entity.Message) error {
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
