package memory

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"time"

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
	conn   *redis.Pool
	keyTTL time.Duration
}

func NewMessageRepository() *messageRepository {
	return &messageRepository{
		conn:   external.RedisPool,
		keyTTL: time.Minute * 10,
	}
}

func (r *messageRepository) IsSingleNoun(message *entity.Message) bool {
	// TODO: +neologd by NewWithDic("path/to/neologd.dic")
	// neologdの辞書のパスの調査
	t := tokenizer.New()
	tokens := t.Tokenize(message.Message)
	if len(tokens) != 3 {
		return false
	}

	firstFeature := tokens[1].Features()

	// "りんご" -> [BOS りんご EOS]
	log.WithFields(log.Fields{
		"message": message.Message,
		"tokens":  tokens,
		"first":   firstFeature,
	}).Info("IsSingleNoun")

	return firstFeature != nil && len(firstFeature) >= 1 && firstFeature[0] == "名詞"
}

// redis message repo の命名規則
// publish message '{"userID": "7141-1414-1414...", "message": "hello"}'
// subscribe messae
// roomID:message になるかも

func (r *messageRepository) Publish(message *entity.Message) error {
	mesInRed := messageInRedis{
		from:    CLIENT,
		Message: message,
	}
	mesBytes, err := json.Marshal(&mesInRed)
	if err != nil {
		return xerrors.Errorf("error in json.Marshal: %w", err)
	}

	conn := r.conn.Get()
	defer conn.Close()

	rep, err := conn.Do("PUBLISH", message.RoomID+":message", mesBytes)
	if err != nil {
		return xerrors.Errorf("error in redis publish: %w", err)
	}

	_, err = conn.Do("EXPIRE", message.RoomID+":message", int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in Publish expire: %w", err)
	}

	log.WithFields(log.Fields{
		"rep":    rep,
		"roomId": message.RoomID,
	}).Info("publish reply")

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
			errCh <- xerrors.Errorf("error in subscribe: %w", err)
		}

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				var message entity.Message

				if err := json.Unmarshal(v.Data, &message); err != nil {
					errCh <- xerrors.Errorf("error in json.Unmarshal: %w", err)
				}
				select {
				case <-ctx.Done():
					log.Info("parent ctx done!")
					return
				default:
					log.WithFields(log.Fields{
						"roomId":  message.RoomID,
						"userId":  message.UserID,
						"message": message.Message,
					}).Info("send message:")
					ch <- &message
				}
			case redis.Subscription:
				log.WithFields(log.Fields{
					"channel": v.Channel,
					"kind":    v.Kind,
					"count":   v.Count,
				}).Info("redis subscription:")
				select {
				case <-ctx.Done():
					log.Info("parent ctx done!")
					return
				default:
					continue
				}
			case error:
				select {
				case <-ctx.Done():
					log.Info("parent ctx done!")
					return
				default:
					errCh <- xerrors.Errorf("error in subscribe: %w", v.Error())
				}
			}
		}
	}()
	return ch, errCh
}
