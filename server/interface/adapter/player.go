package adapter

import (
	"github.com/i-pu/word-war/server/domain/entity"
	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
)

func PbUser2Player(user *pb.User, roomID string) *entity.Player {
	return &entity.Player{
		RoomID: roomID,
		UserID: user.UserId,
	}
}

func Player2PbUser(player *entity.Player) *pb.User {
	return &pb.User{
		UserId: player.UserID,
	}
}

