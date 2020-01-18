package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"time"
)

type roomRepository struct {
	conn   *redis.Pool
	keyTTL time.Duration
	key string
}

func NewRoomRepository() *roomRepository {
	return &roomRepository{
		conn:   external.RedisPool,
		keyTTL: time.Minute * 10,
		key: "roomCandidates",
	}
}

func (r *roomRepository) Lock() error {
	key := r.key + ":lock"
	conn := r.conn.Get()
	// while lockできなかったらloop
	for {
		res, err := conn.Do("SETNX", key, "locking")
		if err != nil {
			return xerrors.Errorf("error in lockRoomCandidate: %w", err)
		}

		_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
		if err != nil {
			return xerrors.Errorf("error in lockRoomCandidate expire: %w", err)
		}

		if res != 0 {
			break
		}
	}
	return nil
}

func (r *roomRepository) Unlock() error {
	key := r.key + ":lock"
	conn := r.conn.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

func (r *roomRepository) GetRoomCandidates() ([]string, error) {
	// lock()
	// defer unlock()
	conn := r.conn.Get()

	value, err := redis.Strings(conn.Do("HKEYS", r.key))
	if err != nil {
		return []string{}, xerrors.Errorf("error in GetRoomCandidates(): %w", err)
	}

	_, err = conn.Do("EXPIRE", r.key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return []string{}, xerrors.Errorf("error in GetRoomCandidates() expire: %w", err)
	}

	return value, nil
}

func (r *roomRepository) AddRoomCandidate(roomID string) error {
	conn := r.conn.Get()

	_, err := conn.Do("HSET", r.key, roomID, 0)
	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()
	if err != nil {
		return xerrors.Errorf("error in AddRoomCandidate HSET %s %s %d: %w", r.key, roomID, 0, err)
	}
	return nil
}

func (r *roomRepository) DeleteRoomCandidate(roomID string) error {
	conn := r.conn.Get()

	_, err := conn.Do("HDEL", r.key, roomID)

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	if err != nil {
		return xerrors.Errorf("error in DeleteRoomCandidate HDEL %s %s %d: %w", r.key, roomID, 0, err)
	}
	return nil
}

// delete <roomID>:**
func (r *roomRepository) DeleteRoom(roomID string) error {
	conn := r.conn.Get()
	// https://blog.morugu.com/entry/2018/01/06/233402
	_, err := conn.Do("EVAL", "return redis.call('DEL', unpack(redis.call('KEYS', ARGV[1])))", 0, roomID + ":*")

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	if err != nil {
		return xerrors.Errorf("error in DeleteRoom HDEL %s %s %d: %w", r.key, roomID, 0, err)
	}
	return nil
}
