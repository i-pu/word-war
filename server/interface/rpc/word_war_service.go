package rpc

import (
	"context"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"
)

type wordWarService struct {
	messageUsecase usecase.MessageUsecase
}

func NewWordWarService(userUsecase usecase.MessageUsecase) *wordWarService {
	return &wordWarService{
		messageUsecase: userUsecase,
	}
}

func (s *wordWarService) Game(in *pb.GameRequest, srv pb.WordWar_GameServer) error {
	for {
		message, err := s.messageUsecase.GetMessage("message")
		if err != nil {
			return err
		}
		counter, err := s.messageUsecase.GetNowCounter()
		if err != nil {
			return err
		}
		// 今は累計メッセージが100を超えたら終了するので
		if counter > 100 {
			// TODO: resultを保存する
			return nil
		}

		res := &pb.GameResponse{
			UserId:  in.UserId,
			Message: message,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
}

func (s *wordWarService) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	return nil, nil
}

func (s *wordWarService) Result(ctx context.Context, in *pb.ResultRequest) (*pb.ResultResponse, error) {
	return nil, nil
}
