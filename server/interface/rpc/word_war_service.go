package rpc

import (
	"context"
	"errors"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"
)

type wordWarService struct {
	messageUsecase usecase.MessageUsecase
	counterUsecase usecase.CounterUsecase
	// 個々にいろんなusecaseついかすればよさそう
}

func NewWordWarService(
	messageUsecase usecase.MessageUsecase,
	counterUsecase usecase.CounterUsecase,
) *wordWarService {
	return &wordWarService{
		messageUsecase: messageUsecase,
		counterUsecase: counterUsecase,
	}
}

func (s *wordWarService) Game(in *pb.GameRequest, srv pb.WordWar_GameServer) error {
	for {
		message, err := s.messageUsecase.GetMessage("message")
		if err != nil {
			return err
		}
		counter, err := s.counterUsecase.Get()
		if err != nil {
			return err
		}
		// 今は累計メッセージが100を超えたら終了するので
		if counter.Value > 100 {
			// TODO: resultを保存する
			return nil
		}

		res := &pb.GameResponse{
			UserId:  message.UserID,
			Message: message.Message,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
}

func (s *wordWarService) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *wordWarService) Result(ctx context.Context, in *pb.ResultRequest) (*pb.ResultResponse, error) {
	return nil, errors.New("not implemented")
}
