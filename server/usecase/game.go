package usecase

import (
	"context"
	"errors"
	"regexp"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type GameUsecase interface {
	InitGameState(roomID string, userID string) error
	TryUpdateWord(message *entity.Message) (*entity.GameState, error)
	SendMessage(message *entity.Message) error
	GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)
	GetCurrentMessage(roomID string) (*entity.Message, error)
}

type gameUsecase struct {
	gameStateRepo  repository.GameStateRepository
	messageRepo    repository.MessageRepository
	messageService *service.MessageService
	counterRepo    repository.CounterRepository
}

func NewGameUsecase(gameRepo repository.GameStateRepository, messageRepo repository.MessageRepository, messageService *service.MessageService, counterRepo repository.CounterRepository) *gameUsecase {
	return &gameUsecase{
		gameStateRepo:  gameRepo,
		messageRepo:    messageRepo,
		messageService: messageService,
		counterRepo:    counterRepo,
	}
}

func (u *gameUsecase) InitGameState(roomID string, userID string) error {
	err := u.gameStateRepo.InitWord(roomID, "しりとり")
	if err != nil {
		return xerrors.Errorf("InitGameState can't InitWord: %w", err)
	}
	if err := u.gameStateRepo.AddUser(roomID, userID); err != nil {
		return xerrors.Errorf("InitGameState can't AddUser(%s, %s): %w", roomID, userID, err)
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
			"TryUpdateWord can't LockCurrentWord. roomId: %v, userId: %v.: %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}
	defer u.gameStateRepo.UnlockCurrentWord(message.RoomID)

	currentWord, err := u.gameStateRepo.GetCurrentWord(message.RoomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't GetCurrentWord. roomId: %v, userId: %v.: %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}

	r := regexp.MustCompile(`^[\p{Hiragana}ー]+$`)
	if !r.Match([]byte(message.Message)) {
		log.WithFields(log.Fields{
			"reason":      "ひらがなでない",
			"currentWord": currentWord,
			"newMessage":  message.Message,
		}).Info("")
		//ひらがなじゃない
		return nil, nil
	}

	if !u.messageRepo.ContainWord(message.Message) {
		log.WithFields(log.Fields{
			"reason":      "存在しない単語",
			"currentWord": currentWord,
			"newMessage":  message.Message,
		}).Info()
		return nil, nil
	}

	if !isSiritori(currentWord, message.Message) {
		log.WithFields(log.Fields{
			"reason":      "しりとりでない",
			"currentWord": currentWord,
			"newMessage":  message.Message,
		}).Info("")

		return nil, nil
	}

	if err := u.gameStateRepo.UpdateCurrentWord(message.RoomID, message.Message); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't UpdateCurrentWord. roomId: %v, userId: %v.: %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}

	if _, err := u.counterRepo.IncrCounter(message.RoomID); err != nil {
		return nil, xerrors.Errorf("TryUpdateWord can't IncrCounter. roomId: %v, userId: %v.: %w",
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
			"SendMessage can't Publish. roomId: %v, userId: %v.: %w",
			message.RoomID,
			message.UserID,
			err,
		)
	}
	return nil
}

// GetMessageChan ctx is used to get cancel signal from parent to cancel pub/sub job
// , so this ctx must be child context.
// repositoryからきたchannelの中身を確認して、
func (u *gameUsecase) GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error) {
	// TODO: 時間制にする
	messageRepoChan, errRepoChan := u.messageRepo.Subscribe(ctx, roomID)

	messageChan := make(chan *entity.Message)
	errChan := make(chan error)
	go func() {
		defer close(messageChan)
		defer close(errChan)
		for {
			select {
			case message, ok := <-messageRepoChan:
				if !ok {
					errChan <- errors.New("logical error about redis channel")
				}

				// ! 10件にしましょう
				counter, err := u.counterRepo.GetCounter(roomID)
				if err != nil {
					errChan <- errors.New("error in counterUsecase.GetScore")
				}
				if counter > 10 {
					// 終了処理
					log.WithFields(log.Fields{
						"roomId": roomID,
					}).Info("Finish game in GetMessageChan.")
					messageChan <- nil
					return
				}
				messageChan <- message

			case err := <-errRepoChan:
				errChan <- xerrors.Errorf("error in game: %w", err)
			}
		}
	}()
	return messageChan, errChan
}

func (u *gameUsecase) GetCurrentMessage(roomID string) (*entity.Message, error) {
	mes, err := u.gameStateRepo.GetCurrentWord(roomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"GetCurrentMessage can't GetCurrentWord. roomId: %v.: %w",
			roomID,
			err,
		)
	}

	return &entity.Message{RoomID: roomID, UserID: "unknown", Message: mes}, nil
}
