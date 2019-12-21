package usecase

import (
	"context"
	"regexp"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
	"golang.org/x/xerrors"
)

type GameUsecase interface {
	InitGameState(roomID string) error
	TryUpdateWord(message *entity.Message) (*entity.GameState, error)
	SendMessage(message *entity.Message) error
	GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)
	GetCurrentMessage(roomID string) (*entity.Message, error)
}

type gameUsecase struct {
	gameStateRepo  repository.GameStateRepository
	messageRepo    repository.MessageRepository
	messageService *service.MessageService
}

func NewMessageUsecase(gameRepo repository.GameStateRepository, messageRepo repository.MessageRepository, messageService *service.MessageService) *gameUsecase {
	return &gameUsecase{
		gameStateRepo:  gameRepo,
		messageRepo:    messageRepo,
		messageService: messageService,
	}
}

func (u *gameUsecase) InitGameState(roomID string) error {
	err := u.gameStateRepo.InitWord(roomID, "しりとり")
	if err != nil {
		return xerrors.Errorf("InitGameState can't InitWord. %w", err)
	}
	return nil
}

// stringは日本語を期待する
func isSiritori(left string, right string) bool {
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	if len(leftRunes) == 0 || len(rightRunes) == 0 {
		return false
	}
	return leftRunes[len(leftRunes)-1] == rightRunes[0]
}

// ひらがな && 1単語 && 名詞
func (u *gameUsecase) TryUpdateWord(message *entity.Message) (*entity.GameState, error) {
	if err := u.gameStateRepo.LockCurrentWord(message.RoomID); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't LockCurrentWord. roomId: %v, userId: %v. %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}
	defer u.gameStateRepo.UnlockCurrentWord(message.RoomID)

	r := regexp.MustCompile(`^\p{Hiragana}+$`)

	currentWord, err := u.gameStateRepo.GetCurrentWord(message.RoomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't GetCurrentWord. roomId: %v, userId: %v. %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}

	if !(isSiritori(currentWord, message.Message) && r.Match([]byte(message.Message)) && u.messageRepo.IsSingleNoun(message)) {
		// 無効なメッセージ
		return nil, nil
	}

	if err := u.gameStateRepo.UpdateCurrentWord(message.RoomID, message.Message); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't UpdateCurrentWord. roomId: %v, userId: %v. %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}

	return &entity.GameState{RoomID: message.RoomID, CurrentWord: message.Message}, nil
}

// SendMessageは周りにメッセージを送る関数
func (u *gameUsecase) SendMessage(message *entity.Message) error {
	if err := u.messageRepo.Publish(message); err != nil {
		return xerrors.Errorf(
			"SendMessage can't Publish. roomId: %v, userId: %v. %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}
	return nil
}

// GetMessageChan ctx is used to get cancel signal from parent to cancel pub/sub job
// , so this ctx must be child context.
func (u *gameUsecase) GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error) {
	messageChan, errChan := u.messageRepo.Subscribe(ctx, roomID)
	return messageChan, errChan
}

func (u *gameUsecase) GetCurrentMessage(roomID string) (*entity.Message, error) {
	mes, err := u.gameStateRepo.GetCurrentWord(roomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"GetCurrentMessage can't GetCurrentWord. roomId: %v. %w",
			roomID,
			err,
		)
	}

	return &entity.Message{RoomID: roomID, UserID: "unknown", Message: mes}, nil
}
