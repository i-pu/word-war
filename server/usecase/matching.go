package usecase

import (
	"fmt"
	"math/rand"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type MatchingUsecase interface {
	Matching(userID string) (string, error)
}

type matchingUsecase struct {
	gameRepo repository.GameRepository
}

func NewMatchingUsecase(gameRepo repository.GameRepository) *matchingUsecase {
	return &matchingUsecase{
		gameRepo: gameRepo,
	}
}

// TODO: すでに部屋が存在した場合はもう一度作成するようにする
// roomIdで10000~99999までの間で5桁の数字を生成
func (u *matchingUsecase) Matching(userID string) (string, error) {
	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("lock matching")

	lockKey := "matching"
	if err := u.gameRepo.Lock(lockKey); err != nil {
		return "", xerrors.Errorf("Lock error: %w", err)
	}
	defer func() {
		log.WithFields(log.Fields{
			"userID": userID,
		}).Debug("unlock")
		if err := u.gameRepo.Unlock(lockKey); err != nil {
			panic(xerrors.Errorf("Unlock error: %w", err))
		}
	}()
	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("GetRoomCandidates")
	rooms, err := u.gameRepo.GetRoomCandidates()
	if err != nil {
		return "", xerrors.Errorf("GetRoomCandidate: %w", err)
	}

	if len(rooms) == 0 {
		roomID := fmt.Sprintf("%d", rand.Intn(90000)+10000)
		player := &entity.Player{RoomID: roomID, UserID: userID}
		log.WithFields(log.Fields{"player": player,}).Debug("len(rooms) == 0")

		if err := u.gameRepo.AddRoomCandidate(roomID); err != nil {
			return "", xerrors.Errorf("AddRoomCandidate(%s) error: %w", roomID, err)
		}
		if err := u.gameRepo.InitWord(roomID, "しりとり"); err != nil {
			return "", xerrors.Errorf("InitWord(%s) error: %w", roomID, err)
		}

		log.WithFields(log.Fields{"player": player}).Debug("AddUser")

		if err := u.gameRepo.AddPlayer(player); err != nil {
			return "", xerrors.Errorf("AddUser(%+v) error: %w", player, err)
		}
		return roomID, nil
	} else {
		roomID := rooms[0]
		player := &entity.Player{RoomID: roomID, UserID: userID}
		if err := u.gameRepo.AddPlayer(player); err != nil {
			return "", xerrors.Errorf("AddUser(%+v) error: %w", player, err)
		}

		userIDs, _ := u.gameRepo.GetUserIDs(roomID)
		// TODO: 待機画面のことを考える。今は空いている部屋があればすぐにroomIDを返すようになっているが、
		// matching rpcをgrpcのstreamで返すようにすれば待機することができるかも
		if len(userIDs) == 4 {
			if err := u.gameRepo.DeleteRoomCandidate(roomID); err != nil {
				return "", xerrors.Errorf("DeleteRoomCandidate(%s) error: %w", roomID, err)
			}
		}

		return roomID, nil
	}
}
