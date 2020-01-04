package main

import (
	"net"

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

// TODO: client用のステータス用のhealthcheckのrcpの作成
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	setUpInfra()
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
	resultUsecase := usecase.NewResultUsecase(resultRepo, resultService)

	pb.RegisterWordWarServer(grpcServer, rpc.NewWordWarService(gameUsecase, counterUsecase, resultUsecase))
	return grpcServer
}

func setUpInfra() {
	external.InitRedis()
}
