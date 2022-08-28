package hallroommgrhandler

import "roomcell/app/hallroommgr/ihallroommgr"

var (
	roomMgrServe ihallroommgr.IHallRoomMgr
)

func InitRoomMgrObj(mgr ihallroommgr.IHallRoomMgr) {
	roomMgrServe = mgr
}
