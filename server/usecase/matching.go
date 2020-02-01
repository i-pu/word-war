package usecase

import (
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type MatchingUsecase interface {
	TryEnterRandomRoom(userID string) (*entity.Room, error)
	CreateRoom(userID string) (*entity.Room, error)
}
type matchingUsecase struct {
	roomRepo repository.RoomRepository
}

func NewMatchingUsecase(roomRepo repository.RoomRepository) *matchingUsecase {
	return &matchingUsecase{
		roomRepo: roomRepo,
	}
}

// TryEnterRandomRoomはrandomなroomに入れなかったら、nilを返す
func (u matchingUsecase) TryEnterRandomRoom(userID string) (*entity.Room, error) {
	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("lock matching")

	lockKey := "matching"
	if err := u.roomRepo.Lock(lockKey); err != nil {
		return nil, xerrors.Errorf("Lock error: %w", err)
	}

	defer func() {
		log.WithFields(log.Fields{
			"userID": userID,
		}).Debug("unlock")
		if err := u.roomRepo.Unlock(lockKey); err != nil {
			panic(xerrors.Errorf("Unlock error: %w", err))
		}
	}()

	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("GetRoomCandidateIDs")

	roomIDs, err := u.roomRepo.GetRoomCandidateIDs()
	if err != nil {
		return nil, xerrors.Errorf("GetRoomCandidate: %w", err)
	}

	if len(roomIDs) == 0 {
		return nil, nil
	} else {
		roomID := roomIDs[0]
		player := &entity.Player{RoomID: roomID, UserID: userID}
		if err := u.roomRepo.AddPlayer(player); err != nil {
			return nil, xerrors.Errorf("AddUser(%+v) error: %w", player, err)
		}

		userIDs, _ := u.roomRepo.GetUserIDs(roomID)
		// TODO: 待機画面のことを考える。今は空いている部屋があればすぐにroomIDを返すようになっているが、
		// matching rpcをgrpcのstreamで返すようにすれば待機することができるかも
		if len(userIDs) == 4 {
			if err := u.roomRepo.DeleteRoomCandidateID(roomID); err != nil {
				return nil, xerrors.Errorf("DeleteRoomCandidateID(%s) error: %w", roomID, err)
			}
		}

		room, err := u.roomRepo.GetRoom(roomID)

		if err != nil {
			return nil, xerrors.Errorf("error in GetRoom(%s) error: %w", roomID, err)
		}

		return room, nil
	}
}

// TODO: RoomIdがかぶった場合などすでに部屋が存在した場合はもう一度作成するようにする
func (u matchingUsecase) CreateRoom(userID string) (*entity.Room, error) {
	room, err := u.roomRepo.CreateRoom()
	if err != nil {
		return nil, xerrors.Errorf("CreateRoom() error: %w", err)
	}

	if err := u.roomRepo.AddRoomCandidateID(room.RoomID); err != nil {
		return nil, xerrors.Errorf("AddRoomCandidateID(%s) error: %w", room.RoomID, err)
	}

	player := &entity.Player{RoomID: room.RoomID, UserID: userID}
	if err := u.roomRepo.AddPlayer(player); err != nil {
		return nil, xerrors.Errorf("AddUser(%+v) error: %w", player, err)
	}

	// workerにenqueueする
	u.roomRepo.StartRoom(room)
	return room, nil
}
