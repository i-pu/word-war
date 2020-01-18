package repository

import "github.com/i-pu/word-war/server/domain/entity"

type ResultRepository interface {
	GetScore(roomID string, userID string) (*entity.Result, error)
	SetScore(roomID string, userID string, score int64) error
	IncrScoreBy(roomID string, userID string, by int64) error
	GetLatestRating(userID string) (int64, error)
	SetRating(userID string, rating int64) error
	AddRatingHistory(userID string, rating int64) error
}
