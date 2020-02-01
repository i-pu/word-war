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
	UpdateRating(player *entity.Player) error
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

func (u *resultUsecase) UpdateRating(player *entity.Player) error {
	// TODO: 部屋に100人いれば100回UpdateRatingが呼ばれるので部屋に固有のgoroutineを作成し、1回だけ呼ばれるようにしたい
	users, err := u.roomRepo.GetUserIDs(player.RoomID)
	if err != nil {
		return xerrors.Errorf("error in UpdateRating(%v). can't roomRepo.GetUsers\n%v", player, err)
	}

	scores := make([]int64, 0, len(users))
	for _, userID := range users {
		p := &entity.Player{RoomID: player.RoomID, UserID: userID}
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
			return xerrors.Errorf("error in UpdateRating(%v). can't roomRepo.GetLatestRating(%s): %w", player, userID, err)
		}
		ratings = append(ratings, rating)
	}

	ratingDeltas, err := ratingDeltasByExtendedEloRating(ratings, sigmoidScores)

	if err != nil {
		return xerrors.Errorf("%v", err)
	}

	for i := 0; i < len(ratingDeltas); i++ {
		log.Debugf("[%d] %d(%d)", i, ratings[i], ratingDeltas[i])
		ratings[i] += ratingDeltas[i]
	}

	for i := 0; i < len(users); i++ {
		if err := u.roomRepo.SetRating(users[i], ratings[i]); err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't roomRepo.SetRating(%s, %d): %w", player.RoomID, users[i], ratings[i], err)
		}
		if err := u.roomRepo.AddRatingHistory(users[i], ratings[i]); err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't roomRepo.AddRatingHistory(%s, %d): %w", player.RoomID, users[i], ratings[i], err)
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
			deltas[i] += int64(elo.RatingDelta(int(ratings[i]), int(ratings[j]), relScoreI))
			relScoreJ := (scores[j] - scores[i] + 1) / 2
			deltas[j] += int64(elo.RatingDelta(int(ratings[i]), int(ratings[j]), relScoreJ))
		}
	}

	return deltas, nil
}
