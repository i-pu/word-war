package repository

import (
	"context"

	"github.com/i-pu/word-war/server/domain/entity"
)

type MessageRepository interface {
	IsSingleNoun(message *entity.Message) bool
	Publish(message *entity.Message) error
	Subscribe(ctx context.Context) (<-chan *entity.Message, <-chan error)
}
