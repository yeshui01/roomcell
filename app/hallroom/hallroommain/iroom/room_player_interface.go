package iroom

import (
	"roomcell/pkg/pb/pbclient"

	"google.golang.org/protobuf/proto"
)

type IGamePlayer interface {
	GetRoleID() int64
	GetName() string
	GetIcon() int32
	SetRoomID(roomID int64)
	GetRoomID() int64
	SendToClient(msgClass int32, msgType int32, pbMsg proto.Message)
	ToClientPlayerInfo() *pbclient.RoomPlayer
}
