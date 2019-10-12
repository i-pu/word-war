package memory

import (
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/infra"
	"log"
	"time"
)

type messageRepository struct {
	conn *redis.Pool
	// roomName  string
	columnKey string
}

func NewMessageRepository() *messageRepository {
	return &messageRepository{
		conn: infra.RedisPool,
		// roomName:  "room1",
	}
}

// redis message repo の命名規則
// publish message '{"userID": "7141-1414-1414...", "message": "hello"}'
// subscribe messae
// 将来 roomID:message になるかも

func (r *messageRepository) Publish(message *entity.Message) error {
	mesBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	conn := r.conn.Get()
	defer conn.Close()
	rep, err := conn.Do("PUBLISH", "message", mesBytes)
	if err != nil {
		return err
	}
	log.Printf("publish reply: %+v", rep)
	return nil
}

// Subscribeのより良いやり方あるかも
func (r *messageRepository) Subscribe(ctx context.Context) (<-chan *entity.Message, error) {
	ch := make(chan *entity.Message)
	go func() {
		defer close(ch)
		conn := r.conn.Get()
		defer conn.Close()
		psc := redis.PubSubConn{Conn: conn}
		_ = psc.Subscribe("message")
		for {
			// 1秒ごとにタイムアウトするのでずっと待ち続けることがなくなる
			switch v := psc.ReceiveWithTimeout(time.Second).(type) {
			case redis.Message:
				var message *entity.Message
				if err := json.Unmarshal(v.Data, message); err != nil {
					panic(err)
				}
				// TODO: ctxが終了したことをチェックしてから送信する
				// こんな適当でいいのだろうか?
				select {
				case <-ctx.Done():
					return
				default:
					ch <- message
				}
			case redis.Subscription:
				log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			case error:
				log.Printf("error: %+v", v.Error())
				// TODO: redisのwithTimeoutのエラーとその他の接続エラーの区別がしたい
				// いま全部ログに出すだけしてるので不安
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}
		}
	}()
	return ch, nil
}
