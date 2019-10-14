package repository

import "github.com/i-pu/word-war/server/domain/entity"

type ResultRepository interface {
	Get(userID string) (*entity.Result, error)
	Set(userID string, result *entity.Result) error
	IncrBy(userID string, by int64) error
}
