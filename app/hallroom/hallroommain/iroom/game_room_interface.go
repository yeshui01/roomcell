package iroom

import (
	"roomcell/pkg/pb/pbclient"

	"google.golang.org/protobuf/proto"
)

type IGameRoom interface {
	GetRoomType() int32
	GetRoomID() int64
	Update(curTime int64)
	JoinPlayer(p IGamePlayer)
	LeavePlayer(roleID int64)
	OnPlayerOffline(p IGamePlayer)
	ToRoomDetail() *pbclient.RoomData
	IsEmpty() bool
	BroadCastRoomMsg(excludeRoleID int64, msgClass int32, msgType int32, pbMsg proto.Message)
	IsCanJoin() bool
	SetMasterID(masterID int64)
	GetMasterID() int64
	IsPlayerFull() bool
}
