package hallroommgr

import (
	"roomcell/app/hallroommgr/hallroommgrhandler"
	"roomcell/app/hallroommgr/hallroommgrmain"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
)

type HallRoomMgr struct {
	mgrGlobalData *hallroommgrmain.RoomMgrGlobal
}

func NewHallRoomMgr() *HallRoomMgr {
	s := &HallRoomMgr{
		mgrGlobalData: hallroommgrmain.NewRommMgr(),
	}
	s.RegisterMsgHandler()
	hallroommgrhandler.InitRoomMgrObj(s)
	return s
}

func (s *HallRoomMgr) GetGlobalData() *hallroommgrmain.RoomMgrGlobal {
	return s.mgrGlobalData
}

func (s *HallRoomMgr) FrameRun(curTimeMs int64) {
	s.mgrGlobalData.Update(curTimeMs)
}

func (s *HallRoomMgr) RegisterMsgHandler() {
	trframe.RegWorkMsgHandler(protocol.ESMsgClassRoom,
		protocol.ESMsgRoomCreate,
		hallroommgrhandler.HandleRoomCreate)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassRoom,
		protocol.ESMsgRoomAutoDelete,
		hallroommgrhandler.HandleRoomAutoDelete)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassRoom,
		protocol.ESMsgRoomFind,
		hallroommgrhandler.HandleRoomFindBrief)

}
