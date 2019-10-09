package rpc

import (
	"context"
	"errors"

	"github.com/i-pu/word-war/server/domain/entity"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"
)

type wordWarService struct {
	// 個々にいろんなusecaseついかすればよさそう
	messageUsecase usecase.MessageUsecase
	counterUsecase usecase.CounterUsecase
	resultUsecase  usecase.ResultUsecase
}

func NewWordWarService(
	messageUsecase usecase.MessageUsecase,
	counterUsecase usecase.CounterUsecase,
	resultUsecase usecase.ResultUsecase,
) *wordWarService {
	return &wordWarService{
		messageUsecase: messageUsecase,
		counterUsecase: counterUsecase,
		resultUsecase:  resultUsecase,
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
		// ! 10件にしましょう
		if counter.Value > 10 {
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
	// 発言者にはそのまま返す
	res := &pb.SayResponse{
		UserId:  in.GetUserId(),
		Message: in.GetMessage(),
	}

	// Redisに(部屋ID固定)
	err := s.messageUsecase.SendMessage("room1", &entity.Message{UserID: in.GetUserId(), Message: in.GetMessage()})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *wordWarService) Result(ctx context.Context, in *pb.ResultRequest) (*pb.ResultResponse, error) {
	return nil, errors.New("not implemented")
}
