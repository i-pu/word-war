package repository

import "github.com/i-pu/word-war/server/domain/entity"

type MessageRepository interface {
	Publish(message *entity.Message) error
	Subscribe(key string) (string, error)
	Set(key string, value string) error
	Get(key string) (string, error)
}
