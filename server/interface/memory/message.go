package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"time"
	"unicode/utf8"

	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
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
	conn       *redis.Pool
	keyTTL     time.Duration
	dictionary *map[string]struct{}
}

func NewMessageRepository() *messageRepository {
	dictionary := map[string]struct{}{}

	if dicPath, ok := os.LookupEnv("DIC_PATH"); !ok {
		log.WithError(xerrors.Errorf("NewMessageRepository LookupEnv(DIC_PATH) is false.")).Warn()
		dictionary["りんご"] = struct{}{}
		dictionary["ごま"] = struct{}{}
		dictionary["まり"] = struct{}{}
	} else {
		fp, err := os.Open(dicPath)
		if err != nil {
			panic(err)
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line := scanner.Text()
			dictionary[line] = struct{}{}
			if !utf8.ValidString(line) {
				log.WithError(xerrors.Errorf("NewMessageRepository utf8.ValidString is false : %w", err)).Fatalf("line: %s", line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.WithError(xerrors.Errorf("NewMessageRepository scanner error: %w", err)).Fatal("")
		}
	}

	return &messageRepository{
		conn:       external.RedisPool,
		keyTTL:     time.Minute * 10,
		dictionary: &dictionary,
	}
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

func (r *messageRepository) ContainWord(word string) bool {
	_, ok := (*r.dictionary)[word]
	return ok
}
