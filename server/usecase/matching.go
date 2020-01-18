package usecase

import (
	"fmt"
	"github.com/i-pu/word-war/server/domain/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"math/rand"
)

type MatchingUsecase interface {
	Matching(userID string) (string, error)
}

type matchingUsecase struct {
	gameStateRepo repository.GameStateRepository
	matchingRepo  repository.RoomRepository
}

func NewMatchingUsecase(gameRepo repository.GameStateRepository, matchingRepo repository.RoomRepository) *matchingUsecase {
	return &matchingUsecase{
		gameStateRepo: gameRepo,
		matchingRepo:  matchingRepo,
	}
}

// TODO: すでに部屋が存在した場合はもう一度作成するようにする
// roomIdで10000~99999までの間で5桁の数字を生成
func (u *matchingUsecase) Matching(userID string) (string, error) {
	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("lock")

	if err := u.matchingRepo.Lock(); err != nil {
		return "", xerrors.Errorf("Lock error: %w", err)
	}
	defer func() {
		log.WithFields(log.Fields{
			"userID": userID,
		}).Debug("unlock")
		if err := u.matchingRepo.Unlock(); err != nil {
			panic(xerrors.Errorf("Unlock error: %w", err))
		}
	}()
	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("GetRoomCandidates")
	rooms, err := u.matchingRepo.GetRoomCandidates()
	if err != nil {
		return "", xerrors.Errorf("GetRoomCandidate: %w", err)
	}

	if len(rooms) == 0 {
		roomID := fmt.Sprintf("%d", rand.Intn(90000)+10000)
		log.WithFields(log.Fields{
			"roomID": roomID,
			"userID": userID,
		}).Debug("len(rooms) == 0")
		if err := u.matchingRepo.AddRoomCandidate(roomID); err != nil {
			return "", xerrors.Errorf("AddRoomCandidate(%s) error: %w", roomID, err)
		}
		if err := u.gameStateRepo.InitWord(roomID, "しりとり"); err != nil {
			return "", xerrors.Errorf("InitWord error: %w", err)
		}

		log.WithFields(log.Fields{
			"roomID": roomID,
			"userID": userID,
		}).Debug("AddUser")
		if err := u.gameStateRepo.AddUser(roomID, userID); err != nil {
			return "", xerrors.Errorf("AddUser error: %w", err)
		}

		return roomID, nil
	} else {
		roomID := rooms[0]
		if err := u.gameStateRepo.AddUser(roomID, userID); err != nil {
			return "", xerrors.Errorf("AddUser error: %w", err)
		}

		users, _ := u.gameStateRepo.GetUsers(roomID)
		// TODO: 待機画面のことを考える。今は空いている部屋があればすぐにroomIDを返すようになっているが、
		// matching rpcをgrpcのstreamで返すようにすれば待機することができるかも
		if len(users) == 4 {
			if err := u.matchingRepo.DeleteRoomCandidate(roomID); err != nil {
				return "", xerrors.Errorf("DeleteRoomCandidate error: %w", err)
			}
		}

		return roomID, nil
	}
}
