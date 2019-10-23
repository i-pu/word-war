package main

import (
	"fmt"
	"log"
	"net"

	"github.com/i-pu/word-war/server/domain/service"
	"github.com/i-pu/word-war/server/infra"
	"github.com/i-pu/word-war/server/interface/memory"
	"github.com/i-pu/word-war/server/interface/rpc"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"

	mecab "github.com/shogo82148/go-mecab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ipadic = "/usr/local/lib/mecab/dic/mecab-ipadic-neologd"
	text   = "dmm.com"
)

func parse(args map[string]string) {
	mecab, err := mecab.New(args)
	if err != nil {
		panic(err)
	}
	defer mecab.Destroy()

	node, err := mecab.ParseToNode(text)

	for ; !node.IsZero(); node = node.Next() {
		fmt.Printf("%s\t%s\n", node.Surface(), node.Feature())
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	setUpInfra()
	grpcServer := setUpGrpc()
	reflection.Register(grpcServer)

	// ====  mecab test  ====
	fmt.Println("# ipadic")
	parse(map[string]string{"dicdir": ipadic})
	// ======================

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
