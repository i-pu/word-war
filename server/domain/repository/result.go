package repository

import "github.com/i-pu/word-war/server/domain/entity"

type ResultRepository interface {
	Get(roomID string) (*entity.Result, error)
	Set(roomID string, result *entity.Result) error
}
