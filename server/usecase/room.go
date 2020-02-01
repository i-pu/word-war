package usecase

import (
	"context"
	"regexp"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type RoomUsecase interface {
	InitUser(player *entity.Player) error
	CleanPlayer(player *entity.Player) error
	TryUpdateWord(message *entity.Message) (*entity.Room, error)
	SendMessage(message *entity.Message) error
	GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)
	GetCurrentMessage(roomID string) (*entity.Message, error)
	GetCounter(roomID string) (int64, error)
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository) *roomUsecase {
	return &roomUsecase{
		roomRepo: roomRepo,
	}
}

// ユーザのゲーム中のデータを削除する。
// 最後のユーザは部屋をきれいにする。複数回読んでも問題ない。resultから呼ばれる
func (u *roomUsecase) CleanPlayer(player *entity.Player) error {
	if err := u.roomRepo.DeletePlayer(player); err != nil {
		return xerrors.Errorf("error roomRepo. CleanPlayer(%+v): %w", player, err)
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

func (u *roomUsecase) InitUser(player *entity.Player) error {
	if err := u.roomRepo.SetScore(player, 0); err != nil {
		return xerrors.Errorf("SetScore(%+v, 0) error: %w", player, err)
	}
	return nil
}

// ひらがな && 1単語 && 名詞
func (u *roomUsecase) TryUpdateWord(message *entity.Message) (*entity.Room, error) {
	lockKey := message.RoomID + ":TryUpdateWord"
	if err := u.roomRepo.Lock(lockKey); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't Lock(%s): message: %+v: %w",
			lockKey,
			message,
			err,
		)
	}
	defer func() {
		if err := u.roomRepo.Unlock(lockKey); err != nil {
			panic(xerrors.Errorf("UnlockCurrentWord(%s): %w", lockKey, err))
		}
	}()

	currentMessage, err := u.roomRepo.GetCurrentMessage(message.RoomID)
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

	if !u.roomRepo.ContainWord(message.Message) {
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

	if err := u.roomRepo.UpdateCurrentMessage(message); err != nil {
		return nil, xerrors.Errorf(
			"TryUpdateWord can't UpdateCurrentWord(%+v): %w",
			message,
			err,
		)
	}

	if _, err := u.roomRepo.IncrCounter(message.RoomID); err != nil {
		return nil, xerrors.Errorf("TryUpdateWord can't IncrCounter(%s): %w",
			message.RoomID,
			err,
		)
	}

	return &entity.Room{RoomID: message.RoomID, CurrentMessage: message}, nil
}

// SendMessageは周りにメッセージを送る関数
func (u *roomUsecase) SendMessage(message *entity.Message) error {
	if err := u.roomRepo.Publish(message); err != nil {
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
func (u *roomUsecase) GetMessageChan(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error) {
	return u.roomRepo.Subscribe(ctx, roomID)
}

func (u *roomUsecase) GetCurrentMessage(roomID string) (*entity.Message, error) {
	mes, err := u.roomRepo.GetCurrentMessage(roomID)
	if err != nil {
		return nil, xerrors.Errorf(
			"GetCurrentMessage can't GetCurrentWord. roomId: %v.: %w",
			roomID,
			err,
		)
	}

	return mes, nil
}

func (u *roomUsecase) GetCounter(roomID string) (int64, error) {
	return u.roomRepo.GetCounter(roomID)
}
