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

func generateMeCab() mecab.MeCab {
	mecab, err := mecab.New(map[string]string{"dicdir": ipadic})
	if err != nil {
		panic(err)
	}
	return mecab
}

func isSingleNoun(mecab *mecab.MeCab, text string) bool {
	node, err := mecab.ParseToNode(text)
	if err != nil {
		panic(err)
	}

	node = node.Next()

	features := strings.Split(node.Feature(), ",")

	fmt.Println(text, features[0], node.Next().Feature(), node.Next().Stat().String() == "EOS")

	return features[0] == "名詞" && node.Next().Stat().String() == "EOS"
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
	mecab := generateMeCab()
	defer mecab.Destroy()
	fmt.Println(isSingleNoun(&mecab, "りんご"))
	fmt.Println(isSingleNoun(&mecab, "動く"))
	fmt.Println(isSingleNoun(&mecab, "青い鳥"))
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
