package ihallroom

import (
	"roomcell/app/hallroom/hallroommain"
)

type IHallRoom interface {
	GetGlobalData() *hallroommain.HallRoomGlobal
	// BroadcastMsgToClient(msgClass int32, msgType int32, pbMsg proto.Message, roleList []int64)
}
