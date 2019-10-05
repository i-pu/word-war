package repository

import "github.com/i-pu/word-war/server/domain/entity"

type MessageRepository interface {
	Publish(key string, message *entity.Message) error
	Subscribe(key string) (string, error)
}
