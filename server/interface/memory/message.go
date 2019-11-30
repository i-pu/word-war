package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	"github.com/ikawaha/kagome/tokenizer"
)

const (
	SERVER = "SERVER"
	CLIENT = "CLIENT"
)

type messageInRedis struct {
	from string `json:"from" validate:"required"`
	*entity.Message
}

type messageRepository struct {
	conn *redis.Pool
}

func NewMessageRepository() *messageRepository {
	return &messageRepository{
		conn: external.RedisPool,
	}
}

func (r *messageRepository) IsSingleNoun(message *entity.Message) bool {
	// [TODO] + neologd by NewWithDic("path/to/neologd.dic")
	t := tokenizer.New()
	tokens := t.Tokenize(message.Message)
	if len(tokens) != 3 {
		return false
	}

	firstFeature := tokens[1].Features()

	// "りんご" -> [BOS りんご EOS]
	fmt.Printf("message: %+v\n tokens: %+v\n first: %+v\n", message, tokens, firstFeature)

	return firstFeature != nil && len(firstFeature) >= 1 && firstFeature[0] == "名詞"
}

// redis message repo の命名規則
// publish message '{"userID": "7141-1414-1414...", "message": "hello"}'
// subscribe messae
// 将来 roomID:message になるかも

func (r *messageRepository) Publish(message *entity.Message) error {
	mesInRed := messageInRedis{
		from:    CLIENT,
		Message: message,
	}
	mesBytes, err := json.Marshal(&mesInRed)
	if err != nil {
		return err
	}
	conn := r.conn.Get()
	defer conn.Close()
	rep, err := conn.Do("PUBLISH", message.RoomID+":message", mesBytes)
	if err != nil {
		return err
	}
	// TODO: loggingの際にroomIDをひょうじしたい
	log.Printf("publish reply: %+v", rep)
	return nil
}

// Subscribeのより良いやり方あるかも
// ctx: 親のcontextで親のcontextが終了するとgo func()内でctx.Done()により終了する
// roomID: どこの部屋のイベントをsubscribeするか
func (r *messageRepository) Subscribe(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error) {
	ch := make(chan *entity.Message)
	errCh := make(chan error)
	go func() {
		defer close(ch)
		defer close(errCh)

		conn := r.conn.Get()
		defer conn.Close()

		psc := redis.PubSubConn{Conn: conn}
		err := psc.Subscribe(roomID + ":message")
		if err != nil {
			errCh <- err
		}

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				var message entity.Message
				log.Printf("%s", string(v.Data))
				if err := json.Unmarshal(v.Data, &message); err != nil {
					errCh <- err
				}
				// TODO: ctxが終了したことをチェックしてから送信する
				// こんな適当でいいのだろうか?
				select {
				case <-ctx.Done():
					log.Printf("parent ctx done!")
					return
				default:
					log.Printf("send message: %v", message)
					ch <- &message
				}
			case redis.Subscription:
				log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
				select {
				case <-ctx.Done():
					log.Printf("parent ctx done!")
					return
				default:
					continue
				}
			case error:
				select {
				case <-ctx.Done():
					log.Printf("parent ctx done!")
					return
				default:
					errCh <- errors.New(v.Error())
				}
			}
		}
	}()
	return ch, errCh
}
