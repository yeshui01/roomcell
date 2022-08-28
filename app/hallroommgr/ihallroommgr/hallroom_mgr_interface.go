package ihallroommgr

import "roomcell/app/hallroommgr/hallroommgrmain"

type IHallRoomMgr interface {
	GetGlobalData() *hallroommgrmain.RoomMgrGlobal
}
