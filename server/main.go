package main

import (
	"fmt"
	"log"
	"net"
	"strings"

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

const ipadic = "/usr/local/lib/mecab/dic/mecab-ipadic-neologd"

func IsSingleNoun(str string) (bool, error) {
	mecab, err := mecab.New(map[string]string{"dicdir": ipadic})
	if err != nil {
		return false, err
	}
	defer mecab.Destroy()

	// XXX: avoid GC problem with MeCab 0.996 (see https://github.com/taku910/mecab/pull/24)
	mecab.Parse("")

	node, err := mecab.ParseToNode(str)

	if err != nil {
		return false, err
	}

	parts := strings.Split(node.Next().Feature(), ",")
	part := parts[0]

	fmt.Println(part)

	// must be a noun
	if part != "名詞" {
		return false, nil
	}

	// must be single part
	if !node.Next().Next().IsZero() {
		return false, nil
	}

	for ; !node.IsZero(); node = node.Next() {
		fmt.Printf("%s \n", strings.Split(node.Feature(), ",")[0])
	}

	return true, nil
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
	fmt.Println(IsSingleNoun("りんご"))
	fmt.Println(IsSingleNoun("動く"))
	fmt.Println(IsSingleNoun("青い鳥"))
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
