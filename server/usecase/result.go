package usecase

import (
	"github.com/google/martian/log"
	"math"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository"
	"github.com/i-pu/word-war/server/domain/service"
	"github.com/kortemy/elo-go"
	"golang.org/x/xerrors"
)

type ResultUsecase interface {
	// userIDはそのuserの結果
	IncrScore(roomID string, userID string, by int64) error
	GetScore(roomID string, userID string) (*entity.Result, error)
	UpdateRating(roomID string, userID string) error
}

type resultUsecase struct {
	resultRepo    repository.ResultRepository
	gameStateRepo repository.GameStateRepository
	service       *service.ResultService
}

func NewResultUsecase(resultRepo repository.ResultRepository, gameStateRepo repository.GameStateRepository, service *service.ResultService) *resultUsecase {
	return &resultUsecase{
		resultRepo:    resultRepo,
		gameStateRepo: gameStateRepo,
		service:       service,
	}
}

func (u *resultUsecase) IncrScore(roomID string, userID string, by int64) error {
	return u.resultRepo.IncrScoreBy(roomID, userID, by)
}

func (u *resultUsecase) GetScore(roomID string, userID string) (*entity.Result, error) {
	return u.resultRepo.GetScore(roomID, userID)
}

func (u *resultUsecase) UpdateRating(roomID string, userID string) error {
	// TODO: 部屋に100人いれば100回UpdateRatingが呼ばれるので部屋に固有のgoroutineを作成し、1回だけ呼ばれるようにしたい
	users, err := u.gameStateRepo.GetUsers(roomID)
	if err != nil {
		return xerrors.Errorf("error in UpdateRating(%s, %s). can't gameStateRepo.GetUsers\n%v", roomID, userID, err)
	}

	scores := make([]int64, 0, len(users))
	for _, user := range users {
		score, err := u.resultRepo.GetScore(roomID, user)
		if err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s).\n%w", roomID, userID, err)
		}
		scores = append(scores, score.Score)
	}

	sigmoidScores := make([]float64, len(scores))
	for i := 0; i < len(scores); i++ {
		sigmoidScores[i] = 1.0 / (1.0 + math.Exp(-float64(scores[i])))
	}

	ratings := make([]int64, 0, len(users))
	for _, user := range users {
		rating, err := u.resultRepo.GetLatestRating(user)
		if err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't resultRepo.GetLatestRating(%s): %w", roomID, userID, user, err)
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
		if err := u.resultRepo.SetRating(users[i], ratings[i]); err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't resultRepo.SetRating(%s, %d): %w", roomID, users[i], ratings[i], err)
		}
		if err := u.resultRepo.AddRatingHistory(users[i], ratings[i]); err != nil {
			return xerrors.Errorf("error in UpdateRating(%s, %s). can't resultRepo.AddRatingHistory(%s, %d): %w", roomID, users[i], ratings[i], err)
		}
	}

	return nil
}

// calculate each member's rating by using extended elo-rating.
// each scores must be in 0 - 1
// rating standard: 1500
func ratingDeltasByExtendedEloRating(ratings []int64, scores []float64) ([]int64, error) {
	if !(len(ratings) == len(scores)) {
		xerrors.New("all args must be same length.")
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
