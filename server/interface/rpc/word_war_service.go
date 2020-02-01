package rpc

import (
	"context"
	"os"
	"strconv"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/repository"
	"github.com/i-pu/word-war/server/usecase"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

type wordWarService struct {
	// 個々にいろんなusecaseついかすればよさそう
	roomUsecase     usecase.RoomUsecase
	resultUsecase   usecase.ResultUsecase
	matchingUsecase usecase.MatchingUsecase
}

func newWordWarService(
	roomUsecase usecase.RoomUsecase,
	resultUsecase usecase.ResultUsecase,
	matchingUsecase usecase.MatchingUsecase,
) *wordWarService {
	return &wordWarService{
		roomUsecase:     roomUsecase,
		resultUsecase:   resultUsecase,
		matchingUsecase: matchingUsecase,
	}
}

func NewGRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer()

	roomRepo := repository.NewRoomRepository()

	roomUsecase := usecase.NewRoomUsecase(roomRepo)
	resultUsecase := usecase.NewResultUsecase(roomRepo)
	matchingUsecase := usecase.NewMatchingUsecase(roomRepo)

	pb.RegisterWordWarServer(grpcServer, newWordWarService(roomUsecase, resultUsecase, matchingUsecase))

	return grpcServer
}

func (s *wordWarService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	serverVersion, versionOk := os.LookupEnv("SERVER_VERSION")
	if !versionOk {
		return nil, xerrors.New("Not Found SERVER_VERSION")
	}
	err := external.HealthCheck()
	redisOk := true
	if err != nil {
		log.Error(xerrors.Errorf("HealthCheck failed: %w", err))
		redisOk = false
	}

	ok := versionOk && redisOk
	ret := &pb.HealthCheckResponse{
		Active:        ok,
		ServerVersion: serverVersion,
	}
	log.WithFields(log.Fields{
		"Active":        ok,
		"ServerVersion": serverVersion,
	}).Info("Health Checked")
	return ret, nil
}

// RoomIDを発行するかもしれないし、すでにあるRoomIDを返すかもしれない
func (s *wordWarService) Matching(ctx context.Context, in *pb.MatchingRequest) (*pb.MatchingResponse, error) {
	// 部屋が存在したら入れる
	room, err := s.matchingUsecase.TryEnterRandomRoom(in.UserId)
	if err != nil {
		return nil, xerrors.Errorf("error: matchingUsecase.TryEnterRandomRoom(%s): %w", err)
	}

	if room == nil {
		room, err = s.matchingUsecase.CreateRoom(in.UserId)
		if err != nil {
			return nil, xerrors.Errorf("error: matchingUsecase.CreateRoom(%s): %w", err)
		}
	}

	ret := &pb.MatchingResponse{RoomId: room.RoomID}
	return ret, nil
}

func (s *wordWarService) Game(in *pb.GameRequest, srv pb.WordWar_GameServer) error {
	err := s.roomUsecase.InitUser(&entity.Player{RoomID: in.RoomId, UserID: in.UserId})
	if err != nil {
		log.WithFields(log.Fields{
			"roomId": in.RoomId,
			"userId": in.UserId,
		}).Fatal(xerrors.Errorf("error in Game(): %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 今の単語を教えてあげる
	mes, err := s.roomUsecase.GetCurrentMessage(in.RoomId)
	if err != nil {
		return xerrors.Errorf("Game rpc can't GetCurrentMessage. roomId: %v, userId: %v.: %w", in.RoomId, in.UserId, err)
	}
	if err := srv.Send(&pb.GameResponse{UserId: mes.UserID, RoomId: mes.RoomID, Message: mes.Message}); err != nil {
		return xerrors.Errorf("Game rpc can't Send. roomId: %v, userId: %v.: %w", in.RoomId, in.UserId, err)
	}

	messageChan, errChan := s.roomUsecase.GetMessageChan(ctx, in.RoomId)
	timerCtx, err := s.roomUsecase.GetTimer(&entity.Room{RoomID: in.RoomId})
	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				log.WithFields(log.Fields{
					"roomId": in.RoomId,
				}).Info("already close messageChan.")
				return nil
			}
			res := &pb.GameResponse{
				UserId:  message.UserID,
				Message: message.Message,
				RoomId:  message.RoomID,
			}
			if err := srv.Send(res); err != nil {
				return xerrors.Errorf("Game rpc can't Send. roomId: %v, userId: %v. : %w", in.RoomId, in.UserId, err)
			}

		case err, ok := <-errChan:
			if !ok {
				log.WithFields(log.Fields{
					"roomId": in.RoomId,
				}).Info("already close errChan.")
				return nil
			}
			log.Error(xerrors.Errorf("error in game rpc: %w", err))

		case <-timerCtx.Done():
			log.WithFields(log.Fields{
				"roomID": in.RoomId,
			}).Info("Finish game in timerCtx.")
			return nil
		}
	}
}

func (s *wordWarService) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	message := &entity.Message{UserID: in.UserId, Message: in.Message, RoomID: in.RoomId}
	game, err := s.roomUsecase.TryUpdateWord(message)
	if err != nil {
		return nil, xerrors.Errorf("Say rpc can't TryUpdateWord. roomId: %v, userId: %v. : %w", in.RoomId, in.UserId, err)
	}
	if game == nil {
		log.WithFields(log.Fields{
			"roomID":  in.RoomId,
			"userId":  in.UserId,
			"message": in.Message,
		}).Info("Can't update word.")
		// なんにも周りに送らない
		res := &pb.SayResponse{Valid: false, UserId: in.UserId, Message: in.Message, RoomId: in.RoomId}
		return res, nil

	}

	// 有効なメッセージしか送らないようになっているから大丈夫なのでまわりに教える
	err = s.roomUsecase.SendMessage(&entity.Message{UserID: in.UserId, Message: in.Message, RoomID: in.RoomId})
	if err != nil {
		return nil, xerrors.Errorf("Say rpc can't SendMessage. roomId: %v, userId: %v. : %w", in.RoomId, in.UserId, err)
	}

	res := &pb.SayResponse{Valid: true, UserId: in.UserId, Message: in.Message, RoomId: in.RoomId}

	err = s.resultUsecase.IncrScore(&entity.Player{RoomID: in.RoomId, UserID: in.UserId}, 5)
	if err != nil {
		return nil, xerrors.Errorf("Say rpc can't IncrScore. roomId: %v, userId: %v. : %w", in.RoomId, in.UserId, err)
	}

	return res, nil
}

func (s *wordWarService) Result(ctx context.Context, in *pb.ResultRequest) (*pb.ResultResponse, error) {
	result, err := s.resultUsecase.GetScore(&entity.Player{RoomID: in.RoomId, UserID: in.UserId})
	if err != nil {
		return nil, xerrors.Errorf("Result rpc can't GetScore. roomId: %v, userId: %v. : %w", in.RoomId, in.UserId, err)
	}

	num := strconv.FormatInt(result.Score, 10)

	res := &pb.ResultResponse{
		UserId: result.UserID,
		Score:  num,
	}

	return res, nil
}
