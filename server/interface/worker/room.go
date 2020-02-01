package worker

import (
	"fmt"
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

	log.Infof("Hello from Worker %+v", room)

	time.Sleep(time.Second * 5)

	log.Infof("GoodBye from Worker %+v", room)

	return nil
}
