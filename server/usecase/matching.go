package usecase

import (
	"fmt"
	"math/rand"

	"github.com/i-pu/word-war/server/domain/repository"
)

type MatchingUsecase interface {
	Matching(userID string) (string, error)
}

type matchingUsecase struct {
	gameStateRepo repository.GameStateRepository
}

func NewMatchingUsecase(gameRepo repository.GameStateRepository) *matchingUsecase {
	return &matchingUsecase{
		gameStateRepo: gameRepo,
	}
}

// TODO: マッチングアルゴリズムを適応する
// TODO: すでに部屋が存在した場合はもう一度作成するようにする
// roomIdで10000~99999までの間で5桁の数字を生成
func (u *matchingUsecase) Matching(userID string) (string, error) {
	roomID := fmt.Sprintf("%d", rand.Intn(90000)+10000)
	// if err := u.gameStateRepo.AddUser(roomID, userID); err != nil {
	// 	return "", xerrors.Errorf("error in Matching userID %s.: %w", userID, err)
	// }
	return roomID, nil
}
