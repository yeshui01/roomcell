package hallroom

import (
	"math/rand"
	"roomcell/app/hallroom/hallroomhandler"
	"roomcell/app/hallroom/hallroommain"
	"time"
)

type HallRoom struct {
	roomGlobal *hallroommain.HallRoomGlobal
}

func NewHallRoom() *HallRoom {
	s := &HallRoom{
		roomGlobal: hallroommain.NewHallRoomGlobal(),
	}
	rand.Seed(time.Now().UnixNano())
	hallroomhandler.InitRoomObj(s)
	s.RegisterMsgHandler()
	return s
}
func (s *HallRoom) GetGlobalData() *hallroommain.HallRoomGlobal {
	return s.roomGlobal
}
func (s *HallRoom) FrameRun(curTimeMs int64) {
	s.roomGlobal.Update(curTimeMs)
}

// func (s *HallRoom) BroadcastMsgToClient(msgClass int32, msgType int32, pbMsg proto.Message, roleList []int64) {
// 	msgData, err := proto.Marshal(pbMsg)
// 	if err != nil {
// 	}
// 	pushMsg := &pbframe.EFrameMsgBroadcastMsgToClientReq{
// 		MsgClass: msgClass,
// 		MsgType:  msgType,
// 		MsgData:  msgData,
// 	}
// 	trframe.BroadcastMessage(protocol.EMsgClassFrame, protocol.EFrameMsgBroadcastMsgToClient, pushMsg, trnode.ETRNodeTypeHallGate)
// }
