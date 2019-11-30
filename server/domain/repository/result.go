package repository

import "github.com/i-pu/word-war/server/domain/entity"

type ResultRepository interface {
	Get(roomID string, userID string) (*entity.Result, error)
	Set(result *entity.Result) error
	IncrBy(roomID string, userID string, by int64) error
}
