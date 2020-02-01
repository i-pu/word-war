package usecase

import (
	"math"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/repository"
	"github.com/kortemy/elo-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type ResultUsecase interface {
	// userIDはそのuserの結果
	IncrScore(player *entity.Player, by int64) error
	GetScore(player *entity.Player) (*entity.Result, error)
	UpdateRating(room *entity.Room) error
}

type resultUsecase struct {
	roomRepo repository.RoomRepository
}

func NewResultUsecase(roomRepo repository.RoomRepository) *resultUsecase {
	return &resultUsecase{
		roomRepo: roomRepo,
	}
}

func (u *resultUsecase) IncrScore(player *entity.Player, by int64) error {
	return u.roomRepo.IncrScoreBy(player, by)
}

func (u *resultUsecase) GetScore(player *entity.Player) (*entity.Result, error) {
	return u.roomRepo.GetScore(player)
}

func (u *resultUsecase) UpdateRating(room *entity.Room) error {
	users, err := u.roomRepo.GetUserIDs(room.RoomID)
	if err != nil {
		return xerrors.Errorf("error in UpdateRating(%v). can't roomRepo.GetUsers(%s)\n%v", room, room.RoomID, err)
	}

	scores := make([]int64, 0, len(users))
	for _, userID := range users {
		p := &entity.Player{RoomID: room.RoomID, UserID: userID}
		score, err := u.roomRepo.GetScore(p)
		if err != nil {
			return xerrors.Errorf("error in UpdateRating(%v).\n%w", p, err)
		}
		scores = append(scores, score.Score)
	}

	sigmoidScores := make([]float64, len(scores))
	for i := 0; i < len(scores); i++ {
		sigmoidScores[i] = 1.0 / (1.0 + math.Exp(-float64(scores[i])))
	}

	ratings := make([]int64, 0, len(users))
	for _, userID := range users {
		rating, err := u.roomRepo.GetLatestRating(userID)
		if err != nil {
			return xerrors.Errorf("error in UpdateRating(%v). can't roomRepo.GetLatestRating(%s): %w", room, userID, err)
		}
		ratings = append(ratings, rating)
	}

	ratingDeltas, err := ratingDeltasByExtendedEloRating(ratings, sigmoidScores)

	if err != nil {
		return xerrors.Errorf("%v", err)
	}

	log.Debug("[Rating Result]")
	for i := 0; i < len(ratingDeltas); i++ {
		log.Debugf("Player %d: %d(%d)", users[i], ratings[i], ratingDeltas[i])
		ratings[i] += ratingDeltas[i]
	}

	for i := 0; i < len(users); i++ {
		if err := u.roomRepo.SetRating(users[i], ratings[i]); err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't roomRepo.SetRating(%s, %d): %w", room.RoomID, users[i], ratings[i], err)
		}
		if err := u.roomRepo.AddRatingHistory(users[i], ratings[i]); err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't roomRepo.AddRatingHistory(%s, %d): %w", room.RoomID, users[i], ratings[i], err)
		}
	}

	return nil
}

// calculate each member's rating by using extended elo-rating.
// each scores must be in 0 - 1
// rating standard: 1500
func ratingDeltasByExtendedEloRating(ratings []int64, scores []float64) ([]int64, error) {
	if !(len(ratings) == len(scores)) {
		return []int64{}, xerrors.New("all args must be same length.")
	}

	elo := elogo.NewElo()
	deltas := make([]int64, len(scores), len(scores))
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			// normalized score
			relScoreI := (scores[i] - scores[j] + 1) / 2
			outcomeI, outcomeJ := elo.Outcome(int(ratings[i]), int(ratings[j]), relScoreI)
			deltas[i] += int64(outcomeI.Delta)
			deltas[j] += int64(outcomeJ.Delta)
		}
	}

	return deltas, nil
}
