package main

import (
	"context"
	"fmt"
	pb "github.com/i-pu/word-war/server/pb"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

type gRPCServer struct{}

func (s *gRPCServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	res := pb.HelloResponse{
		Message: "Hello " + in.Name,
	}
	return &res, nil
}

func (s *gRPCServer) SayHelloManyTimes(in *pb.HelloRequest, srv pb.Greeter_SayHelloManyTimesServer) error {
	for i := 0; i < 5; i++ {
		if err := srv.Send(&pb.HelloResponse{
			Message: "Hello " + in.Name + " Many times: " + fmt.Sprintf("%d", i),
		}); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}
func main() {
	httpServer := http.Server{
		Addr: ":50051",
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &gRPCServer{})

	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		// CORS対策
		grpcweb.WithOriginFunc(func(origin string) bool {
			log.Printf("origin: %v", origin)
			return true
		}),
	)

	httpServer.Handler = http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			// TODO: 公式のコードのコードがとエラーが出るので変更
			// TODO: 調査
			// if wrappedGrpc.IsGrpcWebRequest(req) {
			// log.Printf("well done!!!")
			//	wrappedGrpc.ServeHTTP(resp, req)
			// }
			// Fall back to other servers.
			// http.DefaultServeMux.ServeHTTP(resp, req)
			wrappedGrpc.ServeHTTP(resp, req)
			log.Printf("req: %+v", req)
		},
	)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}
