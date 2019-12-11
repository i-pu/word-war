package rpc

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/i-pu/word-war/server/domain/entity"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/usecase"
)

type wordWarService struct {
	// 個々にいろんなusecaseついかすればよさそう
	gameUsecase    usecase.GameUsecase
	counterUsecase usecase.CounterUsecase
	resultUsecase  usecase.ResultUsecase
}

func NewWordWarService(
	gameUsecase usecase.GameUsecase,
	counterUsecase usecase.CounterUsecase,
	resultUsecase usecase.ResultUsecase,
) *wordWarService {
	return &wordWarService{
		gameUsecase:    gameUsecase,
		counterUsecase: counterUsecase,
		resultUsecase:  resultUsecase,
	}
}

func (s *wordWarService) Matching(ctx context.Context, in *pb.MatchingRequest) (*pb.MatchingResponse, error) {
	// TODO: RoomIdアルゴリズムを適応する
	ret := &pb.MatchingResponse{RoomId: "hogehoge"}
	return ret, nil
}

// TODO: へやに途中から入った人は今の文字がわからないので教えてあげる
// TODO: だめなメッセージも全員に送るようにしてクライアントで処理してもらう
// TODO: 2回同じ単語はだめなので、履歴を保存して検査する
func (s *wordWarService) Game(in *pb.GameRequest, srv pb.WordWar_GameServer) error {
	err := s.gameUsecase.InitGameState(in.RoomId)
	if err != nil {
		log.Printf("init error: %+v", err)
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	// childのcontext荷関数が終了することを教えてあげる
	defer cancel()
	// TODO: 時間制にする
	// TODO: roomに関する情報を削除するゲームが終わったので、resultのあとでもいいかもしれない
	defer s.counterUsecase.Init(in.RoomId, &entity.Counter{Value: 0, RoomID: in.RoomId})
	// redisに終了をpublishする
	// defer redis.publish(done)

	messageChan, errChan := s.gameUsecase.GetMessage(ctx, in.RoomId)
	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				// channelが先に閉じてることはないはずなので
				return errors.New("logical error about redis channel")
			}
			counter, err := s.counterUsecase.Incr(in.RoomId)
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
				RoomId:  message.RoomID,
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
	message := &entity.Message{UserID: in.UserId, Message: in.Message, RoomID: in.RoomId}
	// TODO: validation message

	game, err := s.gameUsecase.TryUpdateWord(message)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	if game == nil {
		log.Printf("can't try to update message: %+v\n", message)
		// なんにも周りに送らない
		res := &pb.SayResponse{Valid: false, UserId: in.UserId, Message: in.Message, RoomId: in.RoomId}
		return res, nil

	}
	// 有効なメッセージしか送らないようになっているので大丈夫なのでまわりに教える
	err = s.gameUsecase.SendMessage(&entity.Message{UserID: in.UserId, Message: in.Message, RoomID: in.RoomId})
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	// FIXME: redisにsendするデータにFrom属性を追加する
	res := &pb.SayResponse{Valid: true, UserId: in.UserId, Message: in.Message, RoomId: in.RoomId}

	// FIXME: てんすう固定
	err = s.resultUsecase.IncrResult(in.RoomId, in.UserId, 5)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *wordWarService) Result(ctx context.Context, in *pb.ResultRequest) (*pb.ResultResponse, error) {
	// 結果を取得する
	result, err := s.resultUsecase.GetResult(in.RoomId, in.UserId)
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
