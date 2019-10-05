package main

import (
	"github.com/i-pu/word-war/server/domain/service"
	"github.com/i-pu/word-war/server/interface/memory"
	"github.com/i-pu/word-war/server/interface/rpc"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func main() {
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
	repo := memory.NewMessageRepository()
	service := service.NewMessageService(repo)
	usecase := usecase.NewMessageUsecase(repo, service)
	pb.RegisterWordWarServer(grpcServer, rpc.NewWordWarService(usecase))
	return grpcServer
}
