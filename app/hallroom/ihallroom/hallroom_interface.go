package ihallroom

import (
	"roomcell/app/hallroom/hallroommain"
)

type IHallRoom interface {
	GetGlobalData() *hallroommain.HallRoomGlobal
	//BroadCastRoomMsg(excludeRoleID int64, msgClass int32, msgType int32, pbMsg proto.Message)
	// BroadcastMsgToClient(msgClass int32, msgType int32, pbMsg proto.Message, roleList []int64)
}
