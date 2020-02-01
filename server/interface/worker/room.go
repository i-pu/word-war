package worker

import (
	"fmt"
	"github.com/i-pu/word-war/server/repository"
	"github.com/i-pu/word-war/server/usecase"
	"golang.org/x/xerrors"
	
	"os"
	"time"

	"encoding/json"
	"github.com/benmanns/goworker"
	"github.com/i-pu/word-war/server/domain/entity"
	log "github.com/sirupsen/logrus"
)

func Worker() {
	settings := goworker.WorkerSettings{
		URI:            "redis://" + os.Getenv("REDIS_URL") + ":6379",
		Connections:    1,
		Queues:         []string{"rooms"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    1,
		Namespace:      "test:",
		IntervalFloat:  1.0,
	}

	goworker.SetSettings(settings)
	goworker.Register("Room", Start)

	log.Debugf("Worker Starting")

	if err := goworker.Work(); err != nil {
		log.Fatalf("Error: %w", err)
	}

	fmt.Println("after work")
}


func Start(queue string, args ...interface{}) error {
	b, _ := json.Marshal(&args[0])
	var room entity.Room
	_ = json.Unmarshal(b, &room)

	// new room usecase
	r := usecase.NewRoomUsecase(repository.NewRoomRepository())

	limit := time.Second * 10
	if err := r.StartGame(&room, limit); err != nil {
		return xerrors.Errorf("StartGame(%+v, %s): %w", room, limit, err)
	}
	log.Debug("timer done")

	resultUsecase := usecase.NewResultUsecase(repository.NewRoomRepository())
	if err := resultUsecase.UpdateRating(&room); err != nil {
		return xerrors.Errorf("UpdateRating(%+v): %w", room, err)
	}
	log.Debug("done updateRating")

	if err := r.EndGame(&room); err != nil {
		return xerrors.Errorf("EndGame(%+v): %w", room, err)
	}
	return nil
}
