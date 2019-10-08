package repository

import "github.com/i-pu/word-war/server/domain/entity"

type MessageRepository interface {
	Publish(roomID string, message *entity.Message) error
	Subscribe(roomID string) (string, error)
}
