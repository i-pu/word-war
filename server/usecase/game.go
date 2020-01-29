package usecase

import (
	"context"
	"regexp"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type GameUsecase interface {
	InitUser(player *entity.Player) error
	CleanGameState(player *entity.Player) error
	TryUpdateWord(message *entity.Message) (*entity.GameState, error)
	SendMessage(message *entity.Message) error
	GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)
	GetCurrentMessage(roomID string) (*entity.Message, error)
	GetCounter(roomID string) (int64, error)
}

type gameUsecase struct {
	gameRepo repository.GameRepository
}

func NewGameUsecase(gameRepo repository.GameRepository) *gameUsecase {
	return &gameUsecase{
		gameRepo: gameRepo,
	}
}

// ユーザのゲーム中のデータを削除する。
// 最後のユーザは部屋をきれいにする。複数回読んでも問題ない。resultから呼ばれる
func (u *gameUsecase) CleanGameState(player *entity.Player) error {
	if err := u.gameRepo.CleanGame(player); err != nil {
		return xerrors.Errorf("error gameRepo.CleanGame(%+v): %w", player, err)
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

func (u *gameUsecase) InitUser(player *entity.Player) error {
	if err := u.gameRepo.SetScore(player, 0); err != nil {
		return xerrors.Errorf("SetScore(%+v, 0) error: %w", player, err)
	}
	return nil
}

// ひらがな && 1単語 && 名詞
func (u *gameUsecase) TryUpdateWord(message *entity.Message) (*entity.GameState, error) {
	lockKey := message.RoomID + ":TryUpdateWord"
	if err := u.gameRepo.Lock(lockKey); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't Lock(%s): message: %+v: %w",
			lockKey,
			message,
			err,
		)
	}
	defer func() {
		if err := u.gameRepo.Unlock(lockKey); err != nil {
			panic(xerrors.Errorf("UnlockCurrentWord(%s): %w", lockKey, err))
		}
	}()

	currentMessage, err := u.gameRepo.GetCurrentMessage(message.RoomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't GetCurrentWord(%s): message: %+v: %w",
			message.RoomID,
			message,
			err,
		)
	}

	// TODO: 伸ばし棒終わったらそのまえの文字を最後の文字とする
	r := regexp.MustCompile(`^[\p{Hiragana}ー]+$`)
	if !r.Match([]byte(message.Message)) {
		log.WithFields(log.Fields{
			"reason":      "ひらがなでない",
			"currentWord": currentMessage.Message,
			"newMessage":  message.Message,
		}).Debug()
		return nil, nil
	}

	if !u.gameRepo.ContainWord(message.Message) {
		log.WithFields(log.Fields{
			"reason":      "存在しない単語",
			"currentWord": currentMessage.Message,
			"newMessage":  message.Message,
		}).Debug()
		return nil, nil
	}

	if !isSiritori(currentMessage.Message, message.Message) {
		log.WithFields(log.Fields{
			"reason":      "しりとりでない",
			"currentWord": currentMessage.Message,
			"newMessage":  message.Message,
		}).Debug()

		return nil, nil
	}

	if err := u.gameRepo.UpdateCurrentMessage(message); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't UpdateCurrentWord(%+v): %w",
			message,
			err,
		)
	}

	if _, err := u.gameRepo.IncrCounter(message.RoomID); err != nil {
		return nil, xerrors.Errorf("TryUpdateWord can't IncrCounter(%s): %w",
			message.RoomID,
			err,
		)
	}

	return &entity.GameState{RoomID: message.RoomID, CurrentWord: message.Message}, nil
}

// SendMessageは周りにメッセージを送る関数
func (u *gameUsecase) SendMessage(message *entity.Message) error {
	if err := u.gameRepo.Publish(message); err != nil {
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
	return u.gameRepo.Subscribe(ctx, roomID)
}

func (u *gameUsecase) GetCurrentMessage(roomID string) (*entity.Message, error) {
	mes, err := u.gameRepo.GetCurrentMessage(roomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"GetCurrentMessage can't GetCurrentWord. roomId: %v.: %w",
			roomID,
			err,
		)
	}

	return mes, nil
}

func (u *gameUsecase) GetCounter(roomID string) (int64, error) {
	return u.gameRepo.GetCounter(roomID)
}
