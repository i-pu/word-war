package main

import (
	"net"
	"os"

	"github.com/benmanns/goworker"
	"github.com/i-pu/word-war/server/external"
	"github.com/i-pu/word-war/server/interface/rpc"
	"github.com/i-pu/word-war/server/interface/worker"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	setUpExternal()

	go worker.Worker()
	defer goworker.Close()

	grpcServer := rpc.NewGRPCServer()
	reflection.Register(grpcServer)
	log.Info("Start server")

	defer grpcServer.Stop()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setUpExternal() {
	external.InitRedis()
	external.InitFirebase()
}

