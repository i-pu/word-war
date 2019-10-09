package main

import (
	"log"
	"net/http"

	"github.com/i-pu/word-war/server/domain/service"
	"github.com/i-pu/word-war/server/infra"
	"github.com/i-pu/word-war/server/interface/memory"
	"github.com/i-pu/word-war/server/interface/rpc"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func main() {
	setUpInfra()
	grpcServer := setUpGrpc()

	// grpc-webを使うためのコード
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			// すべてのホストから許可するので
			return true
		}),
	)

	httpServer := &http.Server{
		Addr: ":50051",
	}
	httpServer.Handler = http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			wrappedGrpc.ServeHTTP(resp, req)
		},
	)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

func setUpGrpc() *grpc.Server {
	grpcServer := grpc.NewServer()
	messageRepo := memory.NewMessageRepository()
	messageService := service.NewMessageService(messageRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, messageService)

	counterRepo := memory.NewCounterRepository()
	counterService := service.NewCounterService(counterRepo)
	counterUsecase := usecase.NewCounterUsecase(counterRepo, counterService)

	resultRepo := memory.NewResultRepository()
	resultService := service.NewResultService(resultRepo)
	resultUsecase := usecase.NewResultUsecase(resultRepo, resultService)

	pb.RegisterWordWarServer(grpcServer, rpc.NewWordWarService(messageUsecase, counterUsecase, resultUsecase))
	return grpcServer
}

func setUpInfra() {
	infra.InitRedis()
}
