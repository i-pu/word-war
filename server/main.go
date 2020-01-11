package main

import (
	"net"
	"os"

	"github.com/i-pu/word-war/server/domain/service"
	"github.com/i-pu/word-war/server/external"
	"github.com/i-pu/word-war/server/interface/memory"
	"github.com/i-pu/word-war/server/interface/rpc"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: OAuthとかの認証を真面目にやってみたい

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

	grpcServer := setUpGrpc()
	reflection.Register(grpcServer)
	log.Info("Start server")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setUpGrpc() *grpc.Server {
	grpcServer := grpc.NewServer()

	messageRepo := memory.NewMessageRepository()
	messageService := service.NewMessageService(messageRepo)

	counterRepo := memory.NewCounterRepository()
	counterService := service.NewCounterService(counterRepo)
	counterUsecase := usecase.NewCounterUsecase(counterRepo, counterService)

	gameRepo := memory.NewGameStateRepository()
	gameUsecase := usecase.NewGameUsecase(gameRepo, messageRepo, messageService, counterRepo)

	resultRepo := memory.NewResultRepository()
	resultService := service.NewResultService(resultRepo)
	resultUsecase := usecase.NewResultUsecase(resultRepo, gameRepo, resultService)

	matchingUsecase := usecase.NewMatchingUsecase(gameRepo)

	pb.RegisterWordWarServer(grpcServer, rpc.NewWordWarService(gameUsecase, counterUsecase, resultUsecase, matchingUsecase))

	return grpcServer
}

func setUpExternal() {
	external.InitRedis()
	external.InitFirebase()
}

// TODO: goleakでgoroutineの数を計測する
