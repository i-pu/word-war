package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

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
	ctx, cancel := context.WithCancel(context.Background())
	// childのcontext荷関数が終了することを教えてあげる
	defer s.counterUsecase.Init(&entity.Counter{Value: 0})
	defer cancel()

	messageChan, errChan := s.messageUsecase.GetMessage(ctx)
	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				// channelが先に閉じてることはないはずなので
				return errors.New("logical error about redis channel")
			}
			counter, err := s.counterUsecase.Incr()
			if err != nil {
				log.Printf("error in incr: %+v", err)
				return err
			}
			// ! 10件にしましょう
			if counter.Value > 10 {
				// 終了処理
				log.Printf("finish game")
				return nil
			}
			res := &pb.GameResponse{
				UserId:  message.UserID,
				Message: message.Message,
			}
			if err := srv.Send(res); err != nil {
				return err
			}
		case err := <-errChan:
			log.Printf("error in game: %+v", err)
		}
	}
}

func (s *wordWarService) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	// mecab test
	isSingleNoun := s.messageUsecase.JudgeMessage(&entity.Message{UserID: in.GetUserId(), Message: in.GetMessage()})
	fmt.Printf("%v is %v", in.GetMessage(), isSingleNoun)

	// 発言者にはそのまま返す
	res := &pb.SayResponse{
		UserId:  in.GetUserId(),
		Message: in.GetMessage(),
	}

	// FIXME: redisにsendするデータにFrom属性を追加する
	// FIXME: 部屋の機能はまだないので、部屋IDはまだ指定しないようにします
	err := s.messageUsecase.SendMessage(&entity.Message{UserID: in.GetUserId(), Message: in.GetMessage()})
	if err != nil {
		return nil, err
	}

	// TODO: スコア計算はどこで？

	// てんすう固定
	err = s.resultUsecase.IncrResult(in.GetUserId(), 5)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *wordWarService) Result(ctx context.Context, in *pb.ResultRequest) (*pb.ResultResponse, error) {
	// 結果を取得する
	result, err := s.resultUsecase.GetResult(in.GetUserId())
	if err != nil {
		return nil, err
	}

	num := strconv.FormatInt(result.Score, 10)

	res := &pb.ResultResponse{
		UserId: result.UserID,
		Score:  num,
	}

	return res, nil
}
