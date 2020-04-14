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
	IsReady(room *entity.Room) (players []*entity.Player, roomUserLimit uint64, timerSeconds uint64, ok bool, err error)
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
// TODO: 4にんまでしか入れない変更できるようにする
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

// TODO: 部屋の設定を受け取るようにして、redisに部屋の上限と時間の制限を設定できるようにする
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

	log.Debug("Enqueue job")
	return room, nil
}

// WatingGameは人が集まるまでblockし続ける関数
// TODO: roomの情報の中に部屋の人数条件を含めたい
// TODO: IsReadyの方で部屋の設定をしないようにしたい
func (u matchingUsecase) IsReady(room *entity.Room) (players []*entity.Player, roomUserLimit uint64, timerSeconds uint64, ok bool, err error) {
	userIDs, err := u.roomRepo.GetUserIDs(room.RoomID)
	if err != nil {
		return nil,
			0,
			0,
			false,
			xerrors.Errorf("AddRoomCandidateID(%s) error: %w", room.RoomID, err)
	}

	players = []*entity.Player{}
	for _, userID := range userIDs {
		players = append(players, u.roomRepo.GetPlayer(room.RoomID, userID))
	}

	roomUserLimit = 2
	timerSeconds = 10
	ok = len(userIDs) == int(roomUserLimit)
	return
}
