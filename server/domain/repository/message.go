package repository

import (
	"context"
	"github.com/i-pu/word-war/server/domain/entity"
)

type MessageRepository interface {
	Publish(message *entity.Message) error
	Subscribe(ctx context.Context) (<-chan *entity.Message, error)
}
